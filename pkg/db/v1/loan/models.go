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
