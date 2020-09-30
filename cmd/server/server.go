package main

import (
	"context"
	"fmt"
	"gin-rest-api-example/internal/account"
	accountDB "gin-rest-api-example/internal/account/database"
	"gin-rest-api-example/internal/article"
	articleDB "gin-rest-api-example/internal/article/database"
	"gin-rest-api-example/internal/config"
	"gin-rest-api-example/internal/database"
	"gin-rest-api-example/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"net/http"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		runApplication()
	},
}

func newServer(lc fx.Lifecycle, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Start to rest api server :%d", cfg.ServerConfig.Port)
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logging.FromContext(ctx).Infof("Stopped rest api server")
			return srv.Shutdown(ctx)
		},
	})
	return r
}

func printAppInfo(cfg *config.Config) {
	logging.DefaultLogger().Infow("application information", "config", cfg)
}

func loadConfig() (*config.Config, error) {
	return config.Load(configFile)
}

func runApplication() {
	// setup application(di + run server)
	app := fx.New(
		fx.Provide(
			// load config
			loadConfig,
			// setup database
			database.NewDatabase,
			// setup account packages
			accountDB.NewAccountDB,
			account.NewAuthMiddleware,
			account.NewHandler,
			// setup article packages
			articleDB.NewArticleDB,
			article.NewHandler,
			// server
			newServer,
		),
		fx.Invoke(
			account.RouteV1,
			article.RouteV1,
			printAppInfo,
		),
	)
	app.Run()
}
