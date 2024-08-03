package errors

import (
	"log"
)

type Error struct {
	ErrName     string `json:"error"`
	Description string `json:"description"`
	Code        int    `json:"errCode"`
}

var ErrorInfo map[string]*Error

const (
	NoDataFound  string = "NoDataFound"
	BadRequest   string = "BadRequest"
	GetDBError   string = "GetDBError"
	AddDBError   string = "AddDBError"
	DelDBError   string = "DelDBError"
	DefaultError string = "DefaultError"
	Unauthorized string = "Unauthorized"
	DbConnError  string = "DbConnError"
	DbError      string = "DBError"
	UnAuthorized string = "UnAuthorized"
)

func ErrorInit() {
	ErrorInfo = make(map[string]*Error)

	ErrorInfo[DbConnError] = &Error{ErrName: DbConnError, Description: "Failed to connect postgres DB:", Code: 1000}
	ErrorInfo[NoDataFound] = &Error{ErrName: NoDataFound, Description: "No data found in system", Code: 1001}
	ErrorInfo[GetDBError] = &Error{ErrName: DbError, Description: "Failed to get data from system", Code: 1002}
	ErrorInfo[AddDBError] = &Error{ErrName: DbError, Description: "Failed to update/insert data into system", Code: 1003}
	ErrorInfo[DelDBError] = &Error{ErrName: DbError, Description: "Failed to Delete the data from system", Code: 1004}
	ErrorInfo[BadRequest] = &Error{ErrName: BadRequest, Description: "Missing or Invalid input arguments", Code: 1005}

	log.Println("ErrorInit successful")
}

func (e *Error) Error() string {
	if e == nil {
		return "UndefinedError"
	}
	return e.ErrName
}

func (e *Error) GetErrorDetails(errMsg string) Error {
	if e == nil {
		return *ErrorInfo[DefaultError]
	}
	err := *e
	if errMsg != "" {
		err.Description = err.Description + " " + errMsg
	}
	return err
}
