package watcher

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yametech/yamecloud/pkg/action/api"
	"github.com/yametech/yamecloud/pkg/uri"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"sync"
)

type watcherServer struct {
	name   string
	parser uri.Parser
	*api.Server
}

func (s *watcherServer) Name() string { return s.name }

func NewWatcherServer(serviceName string, server *api.Server) *watcherServer {
	watcherServer := &watcherServer{
		name:   serviceName,
		parser: uri.NewURIParser(),
		Server: server,
	}
	watcherServer.Group(fmt.Sprintf("/%s", serviceName)).GET("/watch", watcherServer.watch)

	return watcherServer
}

type watcherEvent struct {
	Type   watch.EventType `json:"type"`
	Object runtime.Object  `json:"object"`
	URL    string          `json:"url"`
	Status int             `json:"status"`
}

func (s *watcherServer) watch(g *gin.Context) {
	watcherEventChan := make(chan watcherEvent, 32)
	errors := make(chan error)
	fullURL := g.Request.URL.String()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		errors <- s.startWatch(fullURL, ctx, watcherEventChan)
	}()

	endEvent := watcherEvent{
		Type:   watch.EventType("STREAM_END"),
		URL:    fullURL,
		Status: 410,
	}

	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Writer.CloseNotify():
			cancel()
			g.SSEvent("", endEvent)
			return false

		case err := <-errors:
			if err == nil {
				return true
			}
			g.SSEvent("", endEvent)
			fmt.Printf("[ERROR] watch backend error %s\n", err)
			return false

		case event, ok := <-watcherEventChan:
			if !ok {
				g.SSEvent("", endEvent)
				return false
			}
			g.SSEvent("", event)
		}

		return true
	})
}

func (s *watcherServer) startWatch(url string, ctx context.Context, writeEventChan chan<- watcherEvent) error {
	uris, err := s.parser.ParseWatch(url)
	if err != nil {
		return err
	}

	closes := make([]chan struct{}, 0)
	wg := sync.WaitGroup{}
	wg.Add(len(uris))
	for _, u := range uris {
		eventC, err := s.Watch(u.Namespace, u.Resource, u.ResourceVersion, "")
		if err != nil {
			return err
		}
		_close := make(chan struct{})
		go func() {
			defer wg.Done()
			for {
				select {
				case event, ok := <-eventC:
					if !ok {
						return
					}
					writeEventChan <- watcherEvent{
						Type:   event.Type,
						Object: event.Object,
					}
				case <-_close:
					return
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		for _, _close := range closes {
			_close <- struct{}{}
		}
	}()

	wg.Wait()

	return nil
}
