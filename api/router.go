package api

import (
	"aspire-assignment/pkg/service"

	"github.com/gin-gonic/gin"
)

func getRouter(obj service.ServiceGroupLayer) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check API can be used for the Kubernetes pod health
	router.GET("/health", obj.Health)

	//v1 APIs
	v1Group := router.Group("v1")
	{
		//loan group
		loanGroup := v1Group.Group("loan")
		{
			// loanGroup.POST("", v1.ApplyLoan)
			// loanGroup.PUT("", v1.ApplyLoan)                         //update the loan requested amount
			// loanGroup.DELETE("", v1.ApplyLoan)                      //cancel the loan requested amount
			// loanGroup.PUT("offer", v1.ApplyLoan)                    //pre-approved offers based on monthly salary or bank account balance
			// loanGroup.GET("status", v1.GetLoanStatus)               //approved, rejected, pending amount
			// loanGroup.GET("transactions", v1.GetPaymentTransaction) //transactions against the loan
			// loanGroup.POST("transact", v1.ProcessLoanPayment)       //payments made
			loanGroup.GET("test", obj.GetV1Service().FuncLoanServiceSample)
		}

		//user group
		userGroup := v1Group.Group("user")
		{
			// userGroup.GET("applicaitons", v1.GetPendingLoans) //fetch all applications which are unassigned
			// userGroup.GET("assign", v1.GetPendingLoans)       //assign a loan application to an approver
			// userGroup.POST("update", v1.UpdateLoanStatus)     //update the loan status for assigned applications
			userGroup.GET("test", obj.GetV1Service().FuncUserMgtServiceSample)
		}
	}

	return router
}
