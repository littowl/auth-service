package transport_http

import (
	"auth-service/transport/transport_http/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HttpServer struct {
	logger   *zap.Logger
	httpPort string
}

func NewHttpServer(logger *zap.Logger, httpPort string) *HttpServer {
	return &HttpServer{
		logger:   logger,
		httpPort: httpPort,
	}
}

func (h *HttpServer) StartHTTPServer(handlers *handlers.BaseHandler) {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) { handlers.Ping(c) })

	router.GET("/info", func(c *gin.Context) { handlers.GetUser(c) })

	router.GET("/refresh", func(c *gin.Context) { handlers.Refresh(c) })

	router.POST("/register", func(c *gin.Context) { handlers.Register(c) })

	router.POST("/login", func(c *gin.Context) { handlers.Login(c) })

	router.PATCH("/change_password", func(c *gin.Context) { handlers.ChangePassword(c) })

	h.logger.Info("HTTP server is running on port", zap.String("port", h.httpPort))
	if err := router.Run(h.httpPort); err != nil {
		h.logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

// 2281337
