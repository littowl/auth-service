package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DomainService interface {
	Register(string, string, string) error
	Login(string, string) (LoginResponse, error)
	GetUser(string) (User, error)
	Refresh(string) (LoginResponse, error)
	ChangePassword(string, string) error
}

type BaseHandler struct {
	logger        *zap.Logger
	domainService DomainService
}

type LoginResponse struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	NewPassword string `jso:"new_password"`
}

type User struct {
	ID    int
	Login string
	Hash  string
	Role  string
}

func NewBaseHandler(logger *zap.Logger, domainService DomainService) *BaseHandler {
	return &BaseHandler{
		logger:        logger,
		domainService: domainService,
	}
}

func (h *BaseHandler) Ping(c *gin.Context) {
	h.logger.Info("Ping request received")
	c.String(http.StatusOK, "auth service ok")
}

func (h *BaseHandler) Register(c *gin.Context) {
	h.logger.Info("Registration request received")

	token := c.Request.Header["Authorization"][0]
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to unmarshal request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	err := h.domainService.Register(token, req.Login, req.Password)
	if err != nil {
		h.logger.Error("Failed to register", zap.Error(err))
		if err.Error() == "Invalid token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.String(http.StatusOK, "user was created")
}

func (h *BaseHandler) Login(c *gin.Context) {
	h.logger.Info("Login request received")

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to unmarshal request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	res, err := h.domainService.Login(req.Login, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  res.Access,
		"refresh": res.Refresh,
	})
}

func (h *BaseHandler) GetUser(c *gin.Context) {
	h.logger.Info("Get user request received")

	token := c.Request.Header["Authorization"][0]

	user, err := h.domainService.GetUser(token)
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		if err.Error() == "Invalid token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"login": user.Login,
		"role":  user.Role,
	})
}

func (h *BaseHandler) Refresh(c *gin.Context) {
	h.logger.Info("Refresh request received")

	token := c.Request.Header["Refresh-Token"][0]

	res, err := h.domainService.Refresh(token)
	if err != nil {
		h.logger.Error("Failed to refresh", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access":  res.Access,
		"refresh": res.Refresh,
	})
}

func (h *BaseHandler) ChangePassword(c *gin.Context) {
	h.logger.Info("Change password request received")

	var req ChangePasswordRequest
	token := c.Request.Header["Authorization"][0]
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to unmarshal request body", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unathorized"})
		return
	}

	err := h.domainService.ChangePassword(token, req.NewPassword)
	if err != nil {
		h.logger.Error("Failed to change password", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	c.JSON(http.StatusOK, "password changed")
}
