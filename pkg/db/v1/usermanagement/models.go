package usermanagement

import "database/sql"

type UserDetails struct {
	UserId         sql.NullInt64
	UserName       sql.NullString
	UserPassword   sql.NullString
	UserType       sql.NullString
	Email          sql.NullString
	Mobile         sql.NullString
	MonthlySalary  sql.NullFloat64
	AccountBalance sql.NullFloat64
	CreatedAt      sql.NullTime
	UpdatedAt      sql.NullTime
}
