package loan

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (obj *loanDb) GetUserLoanInstallments(c *gin.Context, userId int64, loanId int64) ([]InstallmentDetails, error) {
	query := `
		select 
			l.id as loan_id,
			l.amount as loan_amount,
			l.status as loan_status,
			i.id as installment_id,
			i.amount_due,
			i.amount_paid,
			i.status as installment_status,
			i.transaction_id,
			i.installment_num,
			i.due_date,
			i.created_at
		FROM
			loan l
		inner join
			installment i
		ON
			i.loan_id = l.id
		where
			l.user_id = ?
			and l.id = ?
		order by i.installment_num;
		`

	rows, err := obj.dbObj.WithContext(c).Raw(query, userId, loanId).Rows()
	if err != nil {
		log.Printf("failed to fetch loans for the user. Error: %s", err.Error())
		return nil, err
	}
	installments := make([]InstallmentDetails, 0)
	for rows.Next() {
		var installment InstallmentDetails
		err := rows.Scan(&installment.LoanId, &installment.LoanAmount, &installment.LoanStatus, &installment.InstallmentId, &installment.AmountDue, &installment.AmountPaid, &installment.Status, &installment.TransactionId, &installment.InstallmentSeq, &installment.DueDate, &installment.CreatedAt)
		if err != nil {
			log.Printf("failed to scan loan. Error:%s", err.Error())
			return nil, err
		}
		installments = append(installments, installment)
	}
	return installments, nil
}

func (obj *loanDb) UpdateInstallment(c *gin.Context, loanId int64, installments []InstallmentDetails, loanClosed bool) error {
	updateQuery := `
		update 
			installment
		set
			amount_paid = ?,
			amount_due = ?,
			status = ?,
			transaction_id = ?
		where
			installment_num = ?
			and loan_id = ?;
	`
	tx := obj.dbObj.Begin()
	for _, installment := range installments {
		updateTx := tx.WithContext(c).Exec(updateQuery, installment.AmountPaid.Float64, installment.AmountDue.Float64, installment.Status.String, installment.TransactionId.String, installment.InstallmentSeq.Int64, loanId)
		if updateTx.Error != nil {
			log.Println("failed to update installment")
			tx.Rollback()
			return updateTx.Error
		}
	}

	updateLoanQuery := `
			update
				loan
			set
				status = 'PAID'
			where 
				id = ?;
	`
	if loanClosed {
		updateTx := tx.WithContext(c).Exec(updateLoanQuery, loanId)
		if updateTx.Error != nil {
			log.Println("failed to update loan closure")
			tx.Rollback()
			return updateTx.Error
		}
	}

	return tx.Commit().Error
}

func (obj *loanDb) UpdateSingleInstallmentPayment(c *gin.Context, loanId int64, installment InstallmentDetails, loanClosed bool) error {
	updateQuery := `
		update 
			installment
		set
			amount_paid = ?,
			amount_due = ?,
			status = ?,
			transaction_id = ?
		where
			installment_num = ?
			and loan_id = ?;
	`
	tx := obj.dbObj.Begin()
	updateTx := tx.WithContext(c).Exec(updateQuery, installment.AmountPaid.Float64, installment.AmountDue.Float64, installment.Status.String, installment.TransactionId.String, installment.InstallmentSeq.Int64, loanId)
	if updateTx.Error != nil {
		log.Println("failed to update installment")
		return updateTx.Error
	}

	updateLoanQuery := `
			update
				loan
			set
				status = 'PAID'
			where 
				id = ?;
	`
	if loanClosed {
		updateTx := tx.WithContext(c).Exec(updateLoanQuery, loanId)
		if updateTx.Error != nil {
			log.Println("failed to update loan closure")
			tx.Rollback()
			return updateTx.Error
		}
	}
	return tx.Commit().Error
}
