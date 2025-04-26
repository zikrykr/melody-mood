package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	appSetup "github.com/melody-mood/cmd/setup"
	"github.com/melody-mood/config"
	"github.com/melody-mood/constants"
	"github.com/melody-mood/middleware"
	"github.com/sirupsen/logrus"

	recommendationRoutes "github.com/melody-mood/internal/recommendations/routes"
	sessionRoutes "github.com/melody-mood/internal/session/routes"
)

// BaseURL base url of api
const BaseURL = "/api/v1"

func StartServer(setupData appSetup.SetupData) {
	conf := config.GetConfig()

	if conf.App.Env == constants.PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	// GIN Init
	router := gin.Default()
	router.UseRawPath = true

	router.Use(middleware.CORSMiddleware())

	//Init Main APP and Route
	initRoute(router, setupData.InternalApp)

	port := config.GetConfig().Http.Port
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := httpServer.ListenAndServe(); err != nil {
			logrus.Error(fmt.Printf("listen: %s\n", err))
		}
	}()
	logrus.Info("webserver started")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logrus.Panic("Server Shutdown:", err)
	}

	logrus.Info("Server exiting")
}

func initRoute(router *gin.Engine, internalAppStruct appSetup.InternalAppStruct) {
	r := router.Group(BaseURL)

	// session
	sessionGroup := r.Group("/session")
	sessionRoutes.Routes.NewRoutes(sessionGroup, internalAppStruct.Handler.SessionHandler)

	// recommendation
	recommendationGroup := r.Group("/recommendations")
	recommendationGroup.Use(middleware.RateLimitMiddleware(internalAppStruct.RedisClient))
	recommendationRoutes.Routes.NewRoutes(recommendationGroup, internalAppStruct.Handler.RecommendationHandler)
}
