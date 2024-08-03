package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *service) Health(c *gin.Context) {
	c.String(http.StatusOK, "I am Healthy")
}
