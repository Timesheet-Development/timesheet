package res

import (
	"fmt"
	"log"
	"net/http"

	"timesheet/commons/validate"

	"github.com/go-chi/render"
)

//ResponseCode signifies either a positive or a negative outcome
type ResponseCode struct {
	//Code signifies either a positive or a negative response code
	Code string `json:"c"`
	//Message a default message for the code
	Message string `json:"m"`
	//HttpStatus that maps to the response code
	HttpStatus int `json:"-"`
}

//Response is a standard web response
type Response struct {
	//Code signifies either a positive or a negative response code
	*ResponseCode
	//Data is any payload
	Data interface{} `json:"d"`
	//Cause for t
	Cause string `json:"v"`
}

//AppError is an application error that communicates error information through
//well defined error codes
type AppError struct {
	*ResponseCode
	Cause error
}

func (e AppError) Error() string {
	return fmt.Sprintf("code [%s] message[%s] status [%d] cause[%s]", e.Code, e.Message, e.HttpStatus, e.Cause.Error())
}

func IsAppError(e error) bool {
	_, ok := e.(AppError)
	return ok
}

func IsAppErrorEquals(e error, code *ResponseCode) bool {
	log.Printf("Checking for apperror [%T]\n", e)
	if ae, ok := e.(*AppError); ok {
		log.Println("Yes, apperror")
		return ae.ResponseCode == code
	}
	log.Println("No, some other error")
	return false
}

//SendError returns a well-formatted standard error response to the browser
func SendError(w http.ResponseWriter, r *http.Request, e error, verbose bool) {

	cause := ""

	if ve, ok := e.(*validate.ValidationError); ok == true {
		log.Printf("Sending validation-error %d\n", http.StatusBadRequest)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, Response{BadRequest, ve, cause})
	} else if ae, ok := e.(AppError); ok == true {
		log.Printf("Sending app-error %d\n", ae.HttpStatus)
		render.Status(r, ae.HttpStatus)
		if verbose {
			if ae.Cause != nil {
				cause = ae.Cause.Error()
			}
		}
		render.JSON(w, r, Response{ae.ResponseCode, nil, cause})
	} else if ae, ok := e.(*AppError); ok == true {
		log.Printf("Sending app-error pointer reference %d\n", ae.HttpStatus)
		render.Status(r, ae.HttpStatus)
		if ae.Cause != nil {
			cause = ae.Cause.Error()
		}
		render.JSON(w, r, Response{ae.ResponseCode, nil, cause})
	} else {
		log.Println("Sending unknown error")
		render.Status(r, http.StatusInternalServerError)
		if ae.Cause != nil {
			cause = ae.Cause.Error()
		}
		render.JSON(w, r, Response{InternalServerError, nil, cause})
	}
}

//SendResponse returns a well-formatted standard error response to the browser
func SendResponse(w http.ResponseWriter, r *http.Request, code *ResponseCode, data interface{}) {
	render.Status(r, code.HttpStatus)
	render.JSON(w, r, Response{code, data, ""})
}

///// Standard Response Codes ////
var InternalServerError = &ResponseCode{"InternalServerError", "Internal Failure. Please retry", http.StatusInternalServerError}
var NoContent = &ResponseCode{"NoContent", "No content", http.StatusNoContent}
var OK = &ResponseCode{"OK", "Your request is completed", http.StatusOK}
var BadRequest = &ResponseCode{"BadRequest", "One or more validation errors occurred.", http.StatusBadRequest}

//// Standard Database Errors
var DatabaseError = &ResponseCode{"DatabaseError", "Internal Failure. Please retry", http.StatusInternalServerError}
var RecordNotFound = &ResponseCode{"RecordNotFound", "Record not found", http.StatusNotFound}
