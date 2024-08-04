package usermanagement

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"aspire-assignment/pkg/auth"
	"aspire-assignment/pkg/db/v1/usermanagement"
	e "aspire-assignment/pkg/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (obj *userMgtService) UserSignup(c *gin.Context) {
	var (
		request  UserSignupRequest
		response UserSignupResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to signup user"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//hash the password
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("unable to hash password. Error:%s", err.Error())
		response.Errors = append(response.Errors, e.ErrorInfo[e.ConversionError].GetErrorDetails("failed to hash the password"))
		response.Message = "failed to signup user"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	hashedPassword := string(hashedPasswordBytes)

	// err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+"LO"))
	// if err != nil {
	// 	fmt.Println("invalid password")
	// } else {
	// 	fmt.Println("valid password")
	// }

	//add the entry into db
	userDetail := usermanagement.UserDetails{
		UserName:       sql.NullString{String: request.UserName, Valid: true},
		UserPassword:   sql.NullString{String: hashedPassword, Valid: true},
		Email:          sql.NullString{String: request.Email, Valid: true},
		UserType:       sql.NullString{String: request.UserType, Valid: true},
		Mobile:         sql.NullString{String: request.Mobile},
		MonthlySalary:  sql.NullFloat64{Float64: request.MonthlySalary, Valid: true},
		AccountBalance: sql.NullFloat64{Float64: request.BankBalance, Valid: true},
	}

	userId, err := obj.dbObj.AddUser(c, userDetail)
	if err != nil {
		log.Printf("failed to add user. Error: %s", err.Error())
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			response.Errors = append(response.Errors, e.ErrorInfo[e.AddDBError].GetErrorDetails("unique username needed"))
		} else {
			response.Errors = append(response.Errors, e.ErrorInfo[e.AddDBError].GetErrorDetails(err.Error()))
		}
		response.Message = "failed to signup user"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Data = &UserSignup{
		UserName: request.UserName,
		UserId:   userId,
	}
	response.Message = "successfully signed up user"
	c.JSON(http.StatusOK, response)
}

func (obj *userMgtService) UserLogin(c *gin.Context) {
	var (
		request  UserLoginRequest
		response UserLoginResponse
	)
	if err := c.BindJSON(&request); err != nil {
		log.Printf("unable to marshal request. Error:%s", err.Error())
		response.Errors = append(response.Errors, *e.ErrorInfo[e.BadRequest])
		response.Message = "failed to login user"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//fetch password hash for the user
	userDetail, err := obj.dbObj.GetUserByUsername(c, request.UserName)
	if err != nil {
		log.Printf("failed to fetch user data. Error: %s", err.Error())
		response.Errors = append(response.Errors, e.ErrorInfo[e.GetDBError].GetErrorDetails(err.Error()))
		response.Message = "failed to login user"
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	if userDetail.UserId.Int64 == 0 {
		//invalid username sent
		log.Println("username not found")
		response.Errors = append(response.Errors, e.ErrorInfo[e.NoDataFound].GetErrorDetails("username not found"))
		response.Message = "failed to login user"
		c.JSON(http.StatusNotFound, response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.UserPassword.String), []byte(request.Password))
	if err != nil || userDetail.UserName.String != request.UserName {
		log.Println("invalid password")
		response.Errors = append(response.Errors, e.ErrorInfo[e.BadRequest].GetErrorDetails("incorrect username/password"))
		response.Message = "failed to login user"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	exp := time.Now().Add(60 * time.Minute)
	payload := auth.Token{
		UserName: userDetail.UserName.String,
		UserId:   userDetail.UserId.Int64,
		UserType: userDetail.UserType.String,
		Exp:      exp,
	}

	token, err := auth.GenerateJWT(payload)
	if err != nil {
		log.Println("failed to generate JWT")
		response.Errors = append(response.Errors, e.ErrorInfo[e.DefaultError].GetErrorDetails("failed to generate JWT token"))
		response.Message = "failed to login user"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Status = true
	response.Data = &UserLogin{
		Token:  token,
		Expiry: exp.Format("2006-01-02 15:04:05"),
	}
	response.Message = "successfully logged in user"
	c.JSON(http.StatusOK, response)
}
