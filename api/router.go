package api

import (
	"aspire-assignment/pkg/auth"
	"aspire-assignment/pkg/service"

	"github.com/gin-gonic/gin"
)

func getRouter(obj service.ServiceGroupLayer) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check API can be used for the Kubernetes pod health
	router.GET("/health", obj.Health)

	//cred APIs
	credGroup := router.Group("cred")
	{
		credGroup.POST("signup", obj.GetV1Service().UserSignup) //signup as customer or admin
		credGroup.POST("login", obj.GetV1Service().UserLogin)   //login for cutomer / admin
	}

	router.Use(auth.AuthMiddleware())

	//v1 APIs
	v1Group := router.Group("v1")
	{
		//loan group
		loanGroup := v1Group.Group("loan")
		{
			loanGroup.POST("", obj.GetV1Service().CreateLoan)                 //create loan for a user id
			loanGroup.PUT("", obj.GetV1Service().ModifyLoan)                  //update the loan requested amount
			loanGroup.DELETE("", obj.GetV1Service().CancelLoan)               //cancel the loan requested amount
			loanGroup.GET("status", obj.GetV1Service().GetLoans)              // fetch loans against user, approved, rejected, pending amount
			loanGroup.GET("installments", obj.GetV1Service().GetInstallments) //transactions against the loan
			loanGroup.POST("repay", obj.GetV1Service().ProcessLoanPayment)    //payments made
			// loanGroup.PUT("offer", v1.ApplyLoan)                    //pre-approved offers based on monthly salary or bank account balance
		}

		//admin group
		adminGroup := v1Group.Group("admin")
		{
			adminGroup.GET("applications", obj.GetV1Service().GetPendingLoans)         //fetch all applications which are unapproved
			adminGroup.POST("update", obj.GetV1Service().ApproveRejectLoanApplication) //update the loan status for assigned applications
			// adminGroup.GET("assign", v1.GetPendingLoans)       //assign a loan application to an approver
		}
	}

	return router
}
