package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/e11it/ra/internal/app/ra"
	loghandler "github.com/e11it/ra/loghandler"
	"github.com/e11it/ra/metrics"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// type config struct {
// 	APPName  string `default:"app name"`
// 	Addr     string `default:":8080"`
// 	LogLevel string `default:""`

// 	Auth auth.Config

// 	ShutdownTimeout uint `default:"5"`
// }

func init() {
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

/* TODO: REMOVE
type Authorizer interface {
	GetMiddleware() gin.HandlerFunc
}

func createAuthRouter(auth_m Authorizer) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(loghandler.Logger(), gin.Recovery())
	// router.Use(helpers.DebugLogger())

	router.Use(auth_m.GetMiddleware())
	router.GET("/auth", func(c *gin.Context) {
		c.String(http.StatusOK, "Auth")
	})
	return router, nil
}*/

func main() {
	gin.SetMode(gin.ReleaseMode)

	monitor := metrics.NewMonitor()
	ra, err := ra.NewRA(getEnv("RA_CONFIG_FILE", "example/config.yml"), monitor)
	if err != nil {
		log.Fatalln(err)
	}
	router := gin.New()

	// Add handler for /metrics and create metrics
	monitor.Use(router)

	router.Use(loghandler.Logger(), gin.Recovery())
	// router.Use(helpers.DebugLogger())

	//router.Use(metrics.Handler())
	// router.GET("/metrics", func(c *gin.Context) {
	// 	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	// })
	router.GET("/auth", ra.GetAuthMiddlerware(), func(c *gin.Context) {
		c.String(http.StatusOK, "Auth")
	})
	router.GET("/reload", func(c *gin.Context) {
		if err := ra.ReloadHandler(); err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}

		c.String(http.StatusOK, "Reload")
	})

	srv := &http.Server{
		Addr:    ra.GetServerAddr(),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s, ok := <-quit
		if !ok {
			break
		}
		switch s {
		case syscall.SIGHUP:
			// TODO: перегрузка конфига
			ra.ReloadHandler()

			/* REMOVE
			updateConfig(Config, cs)
			auth_m.UpdateAuth(&Config.Auth)
			*/
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("shuting down server...")

			// The context is used to inform the server it has 5 seconds to finish
			// the request it is currently handling
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ra.GetShutdownTimeout())*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Println("server forced to shutdown:", err)
			}

			log.Println("server exiting")
			return
		}
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
