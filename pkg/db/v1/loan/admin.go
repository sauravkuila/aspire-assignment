package loan

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (obj *loanDb) GetUnapprovedLoans(c *gin.Context) ([]UnApprovedLoan, error) {
	query := `
		select 
			l.id as loan_id,
			u.user_name,
			l.amount,
			l.tenure,
			l.status,
			l.created_at
		from
			loan l
		inner join
			user_detail u
		on
			l.user_id = u.id
		where
			u.user_type = 'CUSTOMER'
			and l.status = 'PENDING';
	`

	rows, err := obj.dbObj.WithContext(c).Raw(query).Rows()
	if err != nil {
		log.Printf("failed to fetch pending loans. Error: %s", err.Error())
		return nil, err
	}

	loans := make([]UnApprovedLoan, 0)
	for rows.Next() {
		var loan UnApprovedLoan
		err := rows.Scan(&loan.LoanId, &loan.UserName, &loan.Amount, &loan.Installments, &loan.Status, &loan.CreatedAt)
		if err != nil {
			log.Printf("failed to scan loan. Error:%s", err.Error())
			return nil, err
		}
		loans = append(loans, loan)
	}
	return loans, nil
}

func (obj *loanDb) UpdateUnapprovedLoan(c *gin.Context, loanId int64, approved bool) error {
	updateQuery := `
		update 
			loan
		set
			status = ?
		where
			id = ?
		returning id;
	`
	loanStatus := "REJECTED"
	if approved {
		loanStatus = "APPROVED"
	}

	var updatedLoanId sql.NullInt64
	row := obj.dbObj.WithContext(c).Raw(updateQuery, loanStatus, loanId).Row()
	if row.Err() != nil {
		log.Printf("failed to update loan status. Error :%s", row.Err().Error())
		return row.Err()
	}
	if row.Scan(&updatedLoanId) != nil {
		log.Printf("failed to update loan status. Error :%s", row.Err().Error())
		return row.Err()
	}

	return nil
}

func (obj *loanDb) UpdateAndInsertInstallments(c *gin.Context, loanId int64, installmentAmount float64, installment int64) error {
	updateQuery := `
		update 
			loan
		set
			status = 'APPROVED'
		where
			id = ?
		returning id;
	`
	var updatedLoanId sql.NullInt64
	tx := obj.dbObj.Begin()
	updateTx := tx.WithContext(c).Raw(updateQuery, loanId).Scan(&updatedLoanId)
	if updateTx.Error != nil {
		log.Printf("failed to update loan status. Error :%s", updateTx.Error.Error())
		tx.Rollback()
		return updateTx.Error
	}

	if updatedLoanId.Int64 != loanId {
		tx.Rollback()
		return fmt.Errorf("unable to update loan status")
	}

	insertQuery := `
		insert into
			installment(loan_id,amount_due,status,installment_num,due_date)
		values 
	`
	queryFields := make([]string, 0)
	queryValues := make([]interface{}, 0)
	t1 := time.Now()
	for i := 1; i <= int(installment); i++ {
		queryFields = append(queryFields, "(?,?,?,?,?)")
		queryValues = append(queryValues, loanId, installmentAmount, "PENDING", i, t1)
		t1 = t1.Add(24 * 7 * time.Hour)
	}
	insertQuery += strings.Join(queryFields, ",")
	insertTx := tx.WithContext(c).Exec(insertQuery, queryValues...)
	if insertTx.Error != nil {
		log.Printf("failed to insert installments. Error :%s", insertTx.Error.Error())
		tx.Rollback()
		return insertTx.Error
	}
	return tx.Commit().Error
}
