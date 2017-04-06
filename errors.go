package goutils

type GenericError struct {
	ErrorMessage string
}

// USER ERRORS
const (
	CodeInvalidParams = iota + 410
	CodeMissingParam
	CodeBadParam
	CodeMalformedRequest
	CodeInvalidDate
)
const (
	MsgBadParam         = "BAD_PARAMETER {%s}"
	MsgInvalidDate      = "ERROR_DATE_NOT_VALID"
	MsgInvalidParams    = "INVALID_PARAMS"
	MsgMalformedRequest = "MALFORMED_REQUEST"
	MsgMissingParam     = "MISSING_PARAMETER {%s}"
	MsgTransError       = "TRANS_ERROR"
)

// SERVER ERRORS
const (
	CodeDbConnection = iota + 510
	CodeInternalError
)
const (
	MsgDbConnection  = "DB_CONNECTION"
	MsgInternalError = "INTERNAL_ERROR"
)
