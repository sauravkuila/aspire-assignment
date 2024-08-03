package loan

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (obj *loanService) FuncLoanServiceSample(c *gin.Context) {
	c.JSON(http.StatusOK, &gin.H{"status": "working"})
}
