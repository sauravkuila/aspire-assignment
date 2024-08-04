package loan

import "database/sql"

type LoanDetails struct {
	LoanId       sql.NullInt64
	UserId       sql.NullInt64
	Amount       sql.NullFloat64
	Installments sql.NullInt64
	Status       sql.NullString
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
}

type UnApprovedLoan struct {
	LoanId       sql.NullInt64
	UserName     sql.NullString
	Amount       sql.NullFloat64
	Installments sql.NullInt64
	Status       sql.NullString
	CreatedAt    sql.NullTime
}

type InstallmentDetails struct {
	InstallmentId  sql.NullInt64
	LoanId         sql.NullInt64
	LoanAmount     sql.NullFloat64
	LoanStatus     sql.NullString
	AmountDue      sql.NullFloat64
	AmountPaid     sql.NullFloat64
	Status         sql.NullString
	InstallmentSeq sql.NullInt64
	DueDate        sql.NullTime
	TransactionId  sql.NullString
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
}
