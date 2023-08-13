package server

import (
	"context"
	"net/http"

	"cqhttp-client/config"
	"cqhttp-client/pkg/logutil"

	"github.com/hashicorp/go-retryablehttp"
)

type CQServer struct {
	server   *http.Server
	cronTask *cronTask
}

func NewCQServer(c *config.Config) *CQServer {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Logger = nil
	cr := newCronTask(c.CqhttpAddress, c.CrontablesCfg, retryClient.StandardClient())
	return &CQServer{
		server: &http.Server{
			Addr: c.SvcAddress,
			Handler: &Handle{
				client:        retryClient.StandardClient(),
				CqhttpAddress: c.CqhttpAddress,
				keywords:      c.Keywords,
				toID:          c.ToID,
			},
		},
		cronTask: cr,
	}
}

func (c *CQServer) Init() error {
	err := c.cronTask.init()
	if err != nil {
		return err
	}
	return nil
}

func (c *CQServer) Run() {
	logutil.Info("CQServer is starting...")
	go c.server.ListenAndServe()
	c.cronTask.start()
	return
}

func (c *CQServer) Stop() {
	err := c.server.Shutdown(context.Background())
	if err != nil {
		logutil.Error("webserver stopped error", err)
	} else {
		logutil.Info("webserver is stopped")
	}
	c.cronTask.stop()
	return
}
