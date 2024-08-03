package loan

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func (obj *loanDb) CreateLoan(c *gin.Context, userId int64, amount float64, installments int64) (int64, error) {
	query := `
			insert into
				loan(user_id, amount, installments, status)
			values 
				(?,?,?,'PENDING')
			returning 
				id;
			`
	rows, err := obj.dbObj.WithContext(c).Raw(query, userId, amount, installments).Rows()
	if err != nil {
		log.Printf("failed to create a new loan. Error: %s", err.Error())
		return 0, err
	}

	var (
		loanId sql.NullInt64
	)
	for rows.Next() {
		err := rows.Scan(&loanId)
		if err != nil {
			log.Printf("failed to failed inserted loan id. Error: %s", err.Error())
			return 0, err
		}
	}
	return loanId.Int64, nil
}

func (obj *loanDb) ModifyLoan(c *gin.Context, userId int64, loanId int64, amount float64, installments int64) (int64, error) {
	query := `
			update 
				loan
			set
				amount = ?,
				installments = ?
			where
				id = ?
				and user_id = ?
				and status = 'PENDING'
			returning
				id;
			`
	var id sql.NullInt64
	updateTx := obj.dbObj.WithContext(c).Raw(query, amount, installments, loanId, userId).Scan(&id)
	if updateTx.Error != nil {
		log.Printf("failed to modify loan. Error: %s", updateTx.Error.Error())
		return 0, updateTx.Error
	}
	return id.Int64, nil
}

func (obj *loanDb) CancelLoan(c *gin.Context, userId int64, loanId int64) (int64, error) {
	query := `
			update 
				loan
			set
				status = 'CANCELLED'
			where
				id = ?
				and user_id = ?
				and status = 'PENDING'
			returning
				id;
			`
	var id sql.NullInt64
	updateTx := obj.dbObj.WithContext(c).Raw(query, loanId, userId).Scan(&id)
	if updateTx.Error != nil {
		log.Printf("failed to modify loan. Error: %s", updateTx.Error.Error())
		return 0, updateTx.Error
	}
	return id.Int64, nil
}

func (obj *loanDb) GetAllLoansForAgainstUser(c *gin.Context, userId int64) ([]LoanDetails, error) {
	query := `
		select 
			id, amount, installments, status, created_at
		from
			loan
		where
			user_id = ?;
		`

	rows, err := obj.dbObj.WithContext(c).Raw(query, userId).Rows()
	if err != nil {
		log.Printf("failed to fetch loans for the user. Error: %s", err.Error())
		return nil, err
	}
	loans := make([]LoanDetails, 0)
	for rows.Next() {
		var loan LoanDetails
		err := rows.Scan(&loan.LoanId, &loan.Amount, &loan.Installments, &loan.Status, &loan.CreatedAt)
		if err != nil {
			log.Printf("failed to scan loan. Error:%s", err.Error())
			return nil, err
		}
		loans = append(loans, loan)
	}
	return loans, nil
}
