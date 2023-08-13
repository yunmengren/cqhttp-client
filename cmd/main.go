package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"cqhttp-client/config"
	"cqhttp-client/internal/server"
	"cqhttp-client/pkg/logutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		conf string
		cfg  config.Config
	)
	flag.StringVar(&conf, "conf", "config.yml", "set configure file")
	flag.Parse()

	txt, err := ioutil.ReadFile(conf)
	if os.IsNotExist(err) {
		txt, _ = yaml.Marshal(cfg)
		fmt.Printf("build default config: \n%s\n", string(txt))
		os.Exit(-1)
		return
	} else if err != nil {
		fmt.Printf("read conf %s: %v\n", conf, err)
		os.Exit(-2)
		return
	}
	if err = yaml.Unmarshal(txt, &cfg); err != nil {
		fmt.Printf("convert conf %s: err\n", conf, err)
		os.Exit(-3)
		return
	}

	{
		lvl, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		logutil.SetLogFile(cfg.LogFile)
		logutil.SetLogLevel(lvl)
	}

	svc := server.NewCQServer(&cfg)
	err = svc.Init()
	if err != nil {
		logutil.Error("svc init error:", err)
		return
	}
	stopCh := make(chan struct{})
	go func() {
		defer close(stopCh)
		svc.Run()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	s := <-signalChan

	log := logutil.DefaultLogger()
	log.Infof("receive signal: %v", s)
	svc.Stop()
	<-stopCh
}
