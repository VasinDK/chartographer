package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"this_module/config"
	v1 "this_module/internal/controller/http/v1"
	"this_module/internal/pkg/utils"
	"this_module/internal/usecase"
	"this_module/pkg/httpserver"
	"this_module/pkg/logger"
)

/*
todo:
1. Дополнить логер req_id
2. Тестирование unit и интеграционное
3. Комментарии
*/

func Run(cfg *config.Config) {
	l := logger.New(cfg.Logger.Level)
	// pg := 0
	// uc := usecase.New(pg)
	utils := utils.New()
	imageUseCase := usecase.New(utils)

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
