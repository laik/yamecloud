package shell

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/igm/sockjs-go/sockjs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// EndOfTransmission terminal end of
const EndOfTransmission = "\u0004"

// OP      DIRECTION  FIELD(S) USED  DESCRIPTION
type OP uint8

const (
	// BIND    fe->be     sessionID      Id sent back from TerminalResponse
	BIND = iota // 0
	// STDIN   fe->be     Data           Keystrokes/paste buffer
	STDIN // 1
	// STDOUT  be->fe     Data           Output from the process
	STDOUT // 2
	// RESIZE  fe->be     Rows, Cols     New terminal size
	RESIZE // 3
	// TOAST   be->fe     Data           OOB message to be shown to the user
	TOAST // 4
	// INEXIT
	INEXIT // 5
	// OUTEXIT
	OUTEXIT // 6
)

type Type string

const (
	DebugShell  Type = "debug"
	CommonShell Type = "common"
)

// Global import the package init the session manager
var globalSessionManager *sessionManager

func createGlobalSessionManager(clientSet *kubernetes.Clientset, restCfg *rest.Config) {
	if globalSessionManager == nil {
		globalSessionManager = &sessionManager{
			client:   clientSet.CoreV1().RESTClient(),
			restCfg:  restCfg,
			channels: make(map[string]*sessionChannels),
			lock:     sync.RWMutex{},
		}
	}
}

// sessionManager manager external client v2 connect to kubernetes pod terminal session
type sessionManager struct {
	client   rest.Interface
	restCfg  *rest.Config
	channels map[string]*sessionChannels
	lock     sync.RWMutex
}

func (sm *sessionManager) get(sessionId string) (*sessionChannels, bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	v, exists := sm.channels[sessionId]
	return v, exists

}

func (sm *sessionManager) set(sessionId string, channel *sessionChannels) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	sm.channels[sessionId] = channel
}

// close shuts down the SockJS connection and sends the status code and reason to the clientv2
// Can happen if the process exits or if there is an error starting up the process
// For now the status code is unused and reason is shown to the user (unless "")
func (sm *sessionManager) close(sessionID string, status uint32, reason string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	err := sm.channels[sessionID].session.Close(status, reason)
	if err != nil {
		return
	}
	delete(sm.channels, sessionID)
}

// PtyHandler is what remotecommand expects from a pty
type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

// process executed cmd in the container specified in request and connects it up with the  sessionChannels (a session)
func (sm *sessionManager) process(request *attachPodRequest, cmd []string, pty PtyHandler) error {
	base := []string{"/bin/sh", "-c"}
	base = append(base, cmd...)
	req := sm.client.Post().Resource("pods").Name(request.Name).Namespace(request.Namespace).SubResource("exec")

	option := &v1.PodExecOptions{
		Command: base,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if request.Container != "" {
		option.Container = request.Container
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(sm.restCfg, "POST", req.URL())
	if err != nil {
		return err
	}

	return exec.Stream(
		remotecommand.StreamOptions{
			Stdin:             pty,
			Stdout:            pty,
			Stderr:            pty,
			Tty:               true,
			TerminalSizeQueue: pty,
		})
}

// sessionChannels a http connect
// upgrade to websocket session bind a session channel to backend kubernetes API server with SPDY
type sessionChannels struct {
	id       string
	bound    chan struct{}
	session  sockjs.Session
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
	data     chan []byte
}

// message is the messaging protocol between ShellController and TerminalSession.
type message struct {
	Data, SessionID string
	Rows, Cols      uint16
	Width, Height   uint16
	Op              OP
}

// Next impl sizeChan remote command.TerminalSize
func (s *sessionChannels) Next() *remotecommand.TerminalSize {
	select {
	case size := <-s.sizeChan:
		return &size
	case <-s.doneChan:
		return nil
	}
}

// Write impl io.Writer
func (s *sessionChannels) Write(p []byte) (int, error) {
	msg, err := json.Marshal(message{
		Op:   STDOUT,
		Data: string(p),
	})
	if err != nil {
		return 0, err
	}

	if err = s.session.Send(string(msg)); err != nil {
		return 0, err
	}

	return len(p), nil
}

// Toast can be used to send the user any OOB messages
// hterm puts these in the center of the terminal
func (s *sessionChannels) Toast(p string) error {
	msg, err := json.Marshal(message{
		Op:   TOAST,
		Data: p,
	})
	if err != nil {
		return err
	}
	if err = s.session.Send(string(msg)); err != nil {
		return err
	}

	return nil
}

// Read impl io.Reader
func (s *sessionChannels) Read(p []byte) (n int, err error) {
	m, err := s.session.Recv()
	if err != nil {
		return 0, err
	}
	var msg message
	err = json.Unmarshal([]byte(m), &msg)
	if err != nil {
		return copy(p, EndOfTransmission), err
	}

	switch msg.Op {
	case STDIN:
		return copy(p, msg.Data), nil
	case INEXIT: // exit from client v2 event
		return 0, fmt.Errorf("client v2 exit")
	case RESIZE:
		s.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Width,
			Height: msg.Height,
		}
		return 0, nil
	default:
		return copy(p, EndOfTransmission), fmt.Errorf("unknown message type '%d'", msg.Op)
	}
}

func createAttachHandler(path string) http.Handler {
	sockjs.DefaultOptions.JSessionID = sockjs.DefaultJSessionID
	return sockjs.NewHandler(path, sockjs.DefaultOptions, HandleTerminalSession)
}

func HandleTerminalSession(session sockjs.Session) {
	_, _ = session.Recv()
	buf, err := session.Recv()
	if err != nil {
		return
	}

	msg := &message{}
	if err = json.Unmarshal([]byte(buf), &msg); err != nil {
		log.Printf("handleTerminalSession: can't un marshal (%v): %s", buf, err)
		return
	}
	if msg.Op != BIND {
		log.Printf("handleTerminalSession: expected 'bind' message, got: %s", buf)
		return
	}

	count := 0
RETRY:
	terminalSession, exist := globalSessionManager.get(msg.SessionID)
	if !exist {
		time.Sleep(1 * time.Second)
		count++
		if count >= 3 {
			bs, _ := json.Marshal(
				message{
					Op:   OUTEXIT,
					Data: "connection session expired, please close and reconnect",
				})
			if err := session.Send(string(bs)); err != nil {
				log.Printf("send session message to clientv2 error \r\n")
			}
			return
		}
		goto RETRY
	}

	if terminalSession.id == "" {
		log.Printf("handleTerminalSession: can't find session '%s'", msg.SessionID)
		return
	}
	terminalSession.session = session
	globalSessionManager.set(msg.SessionID, terminalSession)
	terminalSession.bound <- struct{}{}
}

// generateTerminalSessionId generates a random session ID string. The format is not really interesting.
// This ID is used to identify the session when the clientv2 opens the SockJS connection.
// Not the same as the SockJS session id! We can't use that as that is generated
// on the client v2 side and we don't have it yet at this point.

func generateTerminalSessionId() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "abc123"
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id)
}

// isValidShell checks if the shell is an allowed one
func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}

func (sm *sessionManager) remoteExecute(method string, url *url.URL, pty PtyHandler, tty bool) error {

	exec, err := remotecommand.NewSPDYExecutor(sm.restCfg, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(
		remotecommand.StreamOptions{
			Stdin:  pty,
			Stdout: pty,
			Stderr: pty,
			Tty:    tty,

			TerminalSizeQueue: pty,
		},
	)
}

func (sm *sessionManager) getContainerIDByName(pod *v1.Pod, containerName string) (string, error) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.Name != containerName {
			continue
		}
		// #52 if a pod is running but not ready(because of readiness probe), we can connect
		if containerStatus.State.Running == nil {
			return "", fmt.Errorf("container [%s] not running", containerName)
		}

		return containerStatus.ContainerID, nil
	}

	// #14 otherwise we should search for running init containers
	for _, initContainerStatus := range pod.Status.InitContainerStatuses {
		if initContainerStatus.Name != containerName {
			continue
		}
		if initContainerStatus.State.Running == nil {
			return "", fmt.Errorf("init container [%s] is not running", containerName)
		}

		return initContainerStatus.ContainerID, nil
	}

	return "", fmt.Errorf("cannot find specified container %s", containerName)
}

func (sm *sessionManager) lanuchDebugPod(request *attachPodRequest, pty PtyHandler) error {
	pod := v1.Pod{}
	err := sm.
		client.
		Get().
		Resource("pods").
		Namespace(request.Namespace).
		Name(request.Name).
		Do(context.Background()).Into(&pod)

	if err != nil {
		return err
	}

	containerID, err := sm.getContainerIDByName(&pod, request.Container)

	if err != nil {
		return err
	}

	// TODO: refactor as kubernetes api style, reuse rbac mechanism of kubernetes
	var targetHost string
	targetHost = pod.Status.HostIP
	//TODO:fix hardcode. should remove const or configMap.will be daynamic
	agentPort := 10027
	uri, err := url.Parse(fmt.Sprintf("http://%s:%d", targetHost, agentPort))
	if err != nil {
		return err
	}
	uri.Path = fmt.Sprintf("/api/v1/debug")
	params := url.Values{}
	//TODO: hardcode,should use front end. Interactive Design probelm
	image := "nicolaka/netshoot:latest"
	params.Add("image", image)
	params.Add("container", containerID)
	params.Add("verbosity", fmt.Sprintf("%v", "0"))
	hstNm, _ := os.Hostname()
	params.Add("hostname", hstNm)

	params.Add("username", "")
	//TODO: should be set false
	params.Add("lxcfsEnabled", "true")
	params.Add("registrySkipTLS", "false")
	params.Add("authStr", "")

	//TODO: support private registry pull image,just like  harbor.
	//var authStr string
	//registrySecret, err := o.CoreClient.Secrets(o.RegistrySecretNamespace).Get(o.RegistrySecretName, v1.GetOptions{})
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		authStr = ""
	//	} else {
	//		return err
	//	}
	//} else {
	//	authStr = string(registrySecret.Data["authStr"])
	//}

	cmd := []string{"/bin/bash", "-c", "export COLUMNS=1000000000", "&&", "export TERM=xterm"}

	commandBytes, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	params.Add("command", string(commandBytes))
	uri.RawQuery = params.Encode()
	return sm.remoteExecute("POST", uri, pty, true)

}

func (sm *sessionManager) launchCommonPod(request *attachPodRequest, session *sessionChannels) error {

	var err error
	validShells := []string{"bash", "sh", "csh"}
	if isValidShell(validShells, request.Shell) {
		cmd := []string{request.Shell}
		err = globalSessionManager.process(request, cmd, session)
	} else {
		// No shell given or it was not valid: try some shells until one succeeds or all fail
		for _, testShell := range validShells {
			cmd := []string{testShell}
			if err = globalSessionManager.process(request, cmd, session); err == nil {
				break
			}
		}
	}
	return err
}

// waitForTerminal is called from pod attach api as a goroutine
// Waits for the SockJS connection to be opened by the client v2 the session to be bound in handleTerminalSession
func waitForTerminal(request *attachPodRequest, sessionId string) {
	session, exist := globalSessionManager.get(sessionId)
	if !exist {
		return
	}
	<-session.bound

	defer close(session.bound)
	var err error
	if request.ShellType == string(DebugShell) {
		err = globalSessionManager.lanuchDebugPod(request, session)
	}

	if request.ShellType == string(CommonShell) {
		err = globalSessionManager.launchCommonPod(request, session)
	}

	if err != nil {
		globalSessionManager.close(sessionId, 2, err.Error())
		return
	}
	globalSessionManager.close(sessionId, 1, "process exited")

}
