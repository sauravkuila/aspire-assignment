package usermanagement

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (obj *userMgtService) FuncUserMgtServiceSample(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"status": "working"})
}
