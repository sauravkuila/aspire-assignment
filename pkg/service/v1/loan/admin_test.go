package loan

import (
	v1 "aspire-assignment/pkg/db/v1"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_loanService_GetPendingLoans(t *testing.T) {
	type fields struct {
		dbObj v1.V1DBLayer
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &loanService{
				dbObj: tt.fields.dbObj,
			}
			obj.GetPendingLoans(tt.args.c)
		})
	}
}

func Test_loanService_ApproveRejectLoanApplication(t *testing.T) {
	type fields struct {
		dbObj v1.V1DBLayer
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := &loanService{
				dbObj: tt.fields.dbObj,
			}
			obj.ApproveRejectLoanApplication(tt.args.c)
		})
	}
}
