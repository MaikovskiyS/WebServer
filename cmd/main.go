package main

import (
	"os"
	"os/signal"
	"stadium/internal/contorller"
	"stadium/internal/repository"
	"stadium/internal/service"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.Info("Create db connection")
	rep, err := repository.NewRepository("root:Wild54323@tcp(127.0.0.1:3306)/world", logger, time.Second*30)
	if err != nil {
		logrus.Error("can't connect to repository")
		// return
	}
	defer rep.Close()

	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		done <- struct{}{}
	}()

	srv := service.NewService(logger, rep)
	tr := contorller.NewController(logger, srv)
	tr.CreatingEndPoints()

	go func() {
		err := tr.Start()
		if err != nil {
			logger.Error(err)
			return
		}
	}()

	<-done

	tr.Stop()
}
