package response

import (
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"net/http"
)

var (
	// ErrTetapTenangTetapSemangat custom error on unexpected error
	ErrTetapTenangTetapSemangat = CustomError{
		Message:  "Tetap Tenang Tetap Semangat",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrUnauthorized = CustomError{
		Message:  "Unauthorized",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrForbidden = CustomError{
		Message:  "Forbidden",
		HTTPCode: http.StatusForbidden,
	}

	ErrNotFound = CustomError{
		Message:  "Record not exist",
		HTTPCode: http.StatusNotFound,
	}

	ErrUnprocessableEntity = CustomError{
		Message:  "Unprocessable Entity",
		HTTPCode: http.StatusUnprocessableEntity,
	}
)

func ProcessValidationError(e *validator.FieldError) map[string]string {
	return map[string]string{"error": ValidationErrorToText(e)}
}

func ValidationErrorToText(e *validator.FieldError) string {
	switch e.Tag {
	case "required":
		return fmt.Sprintf("%s is required", e.Field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field, e.Param)
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field, e.Param)
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field, e.Param)
	}
	return fmt.Sprintf("%s is not valid", e.Field)
}

// SuccessBody holds data for success response
type SuccessBody struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta"`
}

// ErrorBody holds data for error response
type ErrorBody struct {
	Errors []ErrorInfo `json:"errors"`
	Meta   interface{} `json:"meta"`
}

func (e ErrorBody) Error() string {
	errMsg := "response - errors"
	for _, err := range e.Errors {
		errMsg += fmt.Sprintf("\n\t%s", err.Error())
	}
	return errMsg
}

//ErrorResponse error from http.Response
type ErrorResponse struct {
	ErrorBody  ErrorBody
	HTTPStatus int
}

//NewErrorResponse create error response from http response
//Body close is caller responsibility
func NewErrorResponse(resp *http.Response) (ErrorResponse, error) {
	errResp := ErrorResponse{
		HTTPStatus: resp.StatusCode,
		ErrorBody:  ErrorBody{},
	}

	decoder := json.NewDecoder(resp.Body)
	err := decoder.Decode(&errResp.ErrorBody)

	return errResp, err
}

func (e ErrorResponse) Error() string {
	errMsg := e.ErrorBody.Error()
	errMsg += fmt.Sprintf("\nhttp_status: %d", e.HTTPStatus)
	return errMsg
}

// MetaInfo holds meta data
type MetaInfo struct {
	HTTPStatus int `json:"http_status"`
}

// ErrorInfo holds error detail
type ErrorInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Field   string `json:"field,omitempty"`
}

//Error implement error
func (e ErrorInfo) Error() string {
	return fmt.Sprintf(
		"error - msg: %s, code: %d, field: %s",
		e.Message,
		e.Code,
		e.Field,
	)
}

// CustomError holds data for customized error
type CustomError struct {
	Message  string
	Field    string `json:"error"`
	HTTPCode int
}

// Error is a function to convert error to string.
// It exists to satisfy error interface
func (c CustomError) Error() string {
	return c.Message
}

// BuildSuccess is a function to create SuccessBody
func BuildSuccess(data interface{}, message string, meta interface{}) SuccessBody {
	return SuccessBody{
		Data:    data,
		Message: message,
		Meta:    meta,
	}
}

// BuildError is a function to create ErrorBody
func BuildError(errors ...error) ErrorBody {
	if len(errors) == 0 {
		return InternalServerErrorBody()
	}

	errInfos := []ErrorInfo{}

	for _, err := range errors {
		switch errOrig := err.(type) {
		case CustomError:
			return ErrorBody{
				Errors: []ErrorInfo{
					{
						Message: errOrig.Message,
						Field:   errOrig.Field,
					},
				},
				Meta: MetaInfo{
					HTTPStatus: errOrig.HTTPCode,
				},
			}
		case ErrorInfo:
			errInfos = append(errInfos, errOrig)
		case ErrorBody:
			return errOrig
		case ErrorResponse:
			return errOrig.ErrorBody
		default:
			return InternalServerErrorBody()
		}
	}

	return ErrorBody{
		Errors: errInfos,
	}
}

//InternalServerErrorBody for default internal server error
func InternalServerErrorBody() ErrorBody {
	return ErrorBody{
		Errors: []ErrorInfo{
			{
				Message: ErrTetapTenangTetapSemangat.Message,
				Field:   ErrTetapTenangTetapSemangat.Field,
			},
		},
		Meta: MetaInfo{
			HTTPStatus: ErrTetapTenangTetapSemangat.HTTPCode,
		},
	}
}

// Write is a function to write data in json format
func Write(w http.ResponseWriter, result interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(result)
}

// PaginationMetaInfo represent meta info for pagination endpoint
type PaginationMetaInfo struct {
	HTTPStatus int `json:"http_status"`
	Limit      int `json:"limit"`
	Offset     int `json:"offset"`
	Total      int `json:"total"`
}

// NewPaginationMetaInfo to create pagination meta
func NewPaginationMetaInfo(status, limit, offset, total int) PaginationMetaInfo {
	return PaginationMetaInfo{
		HTTPStatus: status,
		Limit:      limit,
		Offset:     offset,
		Total:      total,
	}
}

// OK wraps success responses
func OK(w http.ResponseWriter, data interface{}, message string) {
	successResponse := BuildSuccess(data, message, MetaInfo{HTTPStatus: http.StatusOK})
	Write(w, successResponse, http.StatusOK)
}

// OKWithMeta write response with 2XX http status code with meta
func OKWithMeta(w http.ResponseWriter, data interface{}, msg string, meta interface{}) error {
	sb := BuildSuccess(data, msg, meta)
	return Write(w, sb, http.StatusOK)
}

// Created wrap create response
func Created(w http.ResponseWriter, data interface{}) {
	successResponse := BuildSuccess(data, "Created", MetaInfo{HTTPStatus: http.StatusCreated})
	Write(w, successResponse, http.StatusCreated)
}
