package server

import (
	"cqhttp-client/config"
	"github.com/robfig/cron/v3"
	"net/http"
)

type cronTask struct {
	CqhttpAddress string
	client        *http.Client
	crontables    []*crontable
	cr            *cron.Cron
}

func newCronTask(cqhttpAdress string, cfg []config.CrontableConfig, client *http.Client) *cronTask {
	var crontables []*crontable
	for _, c := range cfg {
		crontables = append(crontables, &crontable{
			toID:    c.ToID,
			message: c.Message,
			cron:    c.Cron,
		})
	}
	return &cronTask{
		CqhttpAddress: cqhttpAdress,
		client:        client,
		crontables:    crontables,
	}
}

type crontable struct {
	toID    int64
	message string
	cron    string
}

func (c *cronTask) init() error {
	cr := cron.New()
	for _, crontab := range c.crontables {
		cron := crontab.cron
		toId := crontab.toID
		message := crontab.message
		_, err := cr.AddFunc(cron, func() {
			c.process(toId, message)
		})
		if err != nil {
			return err
		}
	}
	c.cr = cr
	return nil
}

func (c *cronTask) start() {
	c.cr.Start()
}

func (c *cronTask) stop() {
	c.cr.Stop()
}

func (c *cronTask) process(toID int64, message string) {
	sendMessage(c.CqhttpAddress, toID, "text", message)
}

func (c *crontable) subWeather() {

}
