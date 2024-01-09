package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"this_module/config"
	v1 "this_module/internal/controller/http/v1"
	"this_module/internal/pkg/utils"
	"this_module/internal/repository"
	"this_module/internal/usecase"
	"this_module/pkg/httpserver"
	"this_module/pkg/logger"
)

/*
todo:
1. Дополнить логер req_id
3. Комментарии
*/

func Run(cfg *config.Config) {
	l := logger.New(cfg.Logger.Level)
	utils := utils.New()
	repo, err := repository.New(cfg.FileStorage.Path)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	imageUseCase := usecase.New(utils, repo)

	router, err := v1.NewRouter(l, imageUseCase)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	s := httpserver.New(router, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	select {
	case err = <-s.Notify:
		l.Error(err.Error())
	case <-ctx.Done():
		l.Info("signal Interrupt")
	}

	s.Stop()

	l.Info("Server stoped")
}
