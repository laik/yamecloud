package shell

import (
	"fmt"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
)

type shellServer struct {
	name string
	*api.Server
	// action services
}

func (s *shellServer) Name() string { return s.name }

func NewShellServer(serviceName string, server *api.Server, clientSet *kubernetes.Clientset, cfg *rest.Config) *shellServer {
	shellServer := &shellServer{
		name:   serviceName,
		Server: server,
	}

	createGlobalSessionManager(clientSet, cfg)

	group := shellServer.Group(fmt.Sprintf("/%s", serviceName))
	serveHttp := wrapH(createAttachHandler(fmt.Sprintf("/%s/shell/pod", serviceName)))

	group.Any("/shell/pod/*path", serveHttp)
	group.GET("/attach/namespace/:namespace/pod/:name/container/:container/:shelltype", shellServer.podAttach)

	return shellServer
}

type attachPodRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Container string `json:"container"`
	Shell     string `json:"shell"`
	ShellType string `json:"shellType"`
	Image     string `json:"image"`
}

func wrapH(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers,X-Access-Token,XKey,Authorization")

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *shellServer) podAttach(g *gin.Context) {
	attachPodRequest := &attachPodRequest{
		Namespace: g.Param("namespace"),
		Name:      g.Param("name"),
		Container: g.Param("container"),
		ShellType: g.Param("shelltype"),
		Shell:     g.Query("shell"),
		Image:     g.Query("image"),
	}

	sessionId := generateTerminalSessionId()
	globalSessionManager.set(
		sessionId,
		&sessionChannels{
			id:       sessionId,
			bound:    make(chan struct{}),
			sizeChan: make(chan remotecommand.TerminalSize),
		})

	go waitForTerminal(attachPodRequest, sessionId)

	g.JSON(http.StatusOK, gin.H{"op": BIND, "sessionId": sessionId})
}
