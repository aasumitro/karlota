package internal

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/account"
	"karlota.aasumitro.id/internal/chat"
	"karlota.aasumitro.id/internal/common"
	"karlota.aasumitro.id/internal/notify"
	"karlota.aasumitro.id/internal/utils/cache"
	"karlota.aasumitro.id/internal/utils/http/middleware"
	"karlota.aasumitro.id/internal/utils/http/wrapper"
	"karlota.aasumitro.id/web"
)

func RunApp(
	ctx context.Context,
	infra *config.Infra,
) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(
		ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// init mem cache
	cache.New(cache.NoExpiration, cache.DefaultCleanupInterval)
	// router engine
	routerEngine := infra.GinEngine
	// register public routes
	registerPublicRoutes(routerEngine)
	// register module providers
	registerModuleProvider(routerEngine, infra)
	// server defines parameters for running an HTTP server.
	read, idle := 10, 120
	readHeaderDuration := time.Second * time.Duration(read)
	idleTimeoutDuration := time.Second * time.Duration(idle)
	server := &http.Server{
		Addr:              config.GetServerAddr(),
		Handler:           routerEngine,
		ReadHeaderTimeout: readHeaderDuration,
		IdleTimeout:       idleTimeoutDuration,
		// Disable WriteTimeout for SSE
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	handleShutdown(server, infra)
}

func handleShutdown(
	server *http.Server,
	infra *config.Infra,
) {
	log.Println("shutting down gracefully, press Ctrl+C again to force")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	timeToHandle := 5
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(timeToHandle)*time.Second)
	defer cancel()
	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err)
	}
	// Close database connections
	if err := infra.SQLPool.Close(); err != nil {
		log.Printf("Error close sqlite connection: %v", err)
	}
	// Close RabbitMQ connections
	if err := infra.RabbitMQPublisher.Close(); err != nil {
		log.Printf("Error close rabbitmq connection: %v", err)
	}
	if err := infra.RabbitMQConsumer.Close(); err != nil {
		log.Printf("Error close rabbitmq connection: %v", err)
	}
	// notify user of shutdown
	log.Println("Server exiting")
}

type embeddedFile struct {
	fs.File
}

func (f *embeddedFile) Close() error {
	return nil
}

func (f *embeddedFile) Seek(
	offset int64,
	whence int,
) (int64, error) {
	return f.File.(io.Seeker).Seek(offset, whence)
}

func registerPublicRoutes(router *gin.Engine) {
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/ui")
	})
	router.StaticFS("/ui", http.FS(web.SPA()))
	router.NoRoute(func(ctx *gin.Context) {
		if !strings.Contains(ctx.FullPath(), "/ui/") {
			ctx.String(http.StatusNotFound,
				"route that you are looking for is not found")
			return
		}
		file, err := web.SPA().Open("index.html")
		if err != nil {
			ctx.String(http.StatusInternalServerError,
				"failed to open spa file: ", err.Error())
			return
		}
		defer func() { _ = file.Close() }()
		fileInfo, err := file.Stat()
		if err != nil {
			ctx.String(http.StatusInternalServerError,
				"failed to get spa file info: ", err.Error())
			return
		}
		http.ServeContent(
			ctx.Writer, ctx.Request, fileInfo.Name(),
			fileInfo.ModTime(), &embeddedFile{file})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		wrapper.NewHTTPRespondWrapper(
			ctx, http.StatusOK, "PONG")
	})
	router.GET("/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(
				common.SwaggerDefaultModelsExpandDepth)))
}

func registerModuleProvider(
	router *gin.Engine,
	infra *config.Infra,
) {
	v1 := router.Group("api/v1")
	middleware.SetRequirement(config.GetAuthSecret())
	account.NewAccountProvider(v1, infra)
	chat.NewChatProvider(v1, infra)
	notify.NewNotificationProvider(infra)
}
