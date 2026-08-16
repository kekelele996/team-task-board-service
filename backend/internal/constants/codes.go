package constants

// Unified API response codes.
const (
	CodeOK           = 0
	CodeBadRequest   = 40000
	CodeUnauthorized = 40100
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeConflict     = 40900
	CodeInternal     = 50000
)

// MessageByCode maps codes to human readable messages.
var MessageByCode = map[int]string{
	CodeOK:           "ok",
	CodeBadRequest:   "bad request",
	CodeUnauthorized: "unauthorized",
	CodeForbidden:    "forbidden",
	CodeNotFound:     "not found",
	CodeConflict:     "conflict",
	CodeInternal:     "internal server error",
}
