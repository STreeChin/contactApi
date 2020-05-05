package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/STreeChin/contactapi/internal/controller"
	"github.com/STreeChin/contactapi/internal/service"
	"github.com/STreeChin/contactapi/pkg/cache"
	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/STreeChin/contactapi/pkg/log"
	"github.com/STreeChin/contactapi/pkg/repository"
	"github.com/STreeChin/contactapi/pkg/route"
	"github.com/STreeChin/contactapi/pkg/route/middleware"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	logger := log.NewLogger(cfg)
	rep := repository.NewRepository(logger, cfg)
	cacher := cache.NewCache(logger, cfg)
	srv := service.NewContactService(logger, cfg, cacher, rep)
	ctrl := controller.NewContactController(logger, srv)
	rtr := route.NewRouter(ctrl)

	authMdw := middleware.NewAuth(logger, rep)
	rtr.Use(authMdw.Middleware)

	port := cfg.Host.Port
	if "" == port {
		port = ":8080"
	}
	fmt.Println("Server started")
	logger.Errorln(http.ListenAndServe(port, rtr))
}
