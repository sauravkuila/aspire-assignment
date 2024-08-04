package usermanagement

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func (obj *userMgtDb) AddUser(c *gin.Context, userDetail UserDetails) (int64, error) {
	query := `
		insert into
			user_detail(user_name,password,user_type,email,mobile,monthly_salary,acc_bal)
		values
			(?,?,?,?,?,?,?)
		returning id;
	`

	var userId sql.NullInt64
	insertTx := obj.dbObj.WithContext(c).Raw(query, userDetail.UserName.String, userDetail.UserPassword.String, userDetail.UserType.String, userDetail.Email.String, userDetail.Mobile.String, userDetail.MonthlySalary.Float64, userDetail.AccountBalance.Float64).Scan(&userId)
	if insertTx.Error != nil {
		log.Println("error in adding user")
		return 0, insertTx.Error
	}

	return userId.Int64, nil
}

func (obj *userMgtDb) GetUserByUsername(c *gin.Context, userName string) (UserDetails, error) {
	query := `
		select 
			id, 
			user_name, 
			password, 
			user_type, 
			email, 
			mobile, 
			monthly_salary, 
			acc_bal, 
			created_at
		from
			user_detail
		where
			user_name=?;
	`

	var userDetail UserDetails
	rows, err := obj.dbObj.WithContext(c).Raw(query, userName).Rows()
	if err != nil {
		log.Println("failed to fetch user detail")
		return userDetail, err
	}
	for rows.Next() {
		err := rows.Scan(&userDetail.UserId, &userDetail.UserName, &userDetail.UserPassword, &userDetail.UserType, &userDetail.Email, &userDetail.Mobile, &userDetail.MonthlySalary, &userDetail.AccountBalance, &userDetail.CreatedAt)
		if err != nil {
			log.Println("failed to scan user detail")
			return userDetail, err
		}
	}
	return userDetail, nil
}
