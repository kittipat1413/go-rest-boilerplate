package constants

const (
	// Success code
	StatusCodeSuccess = "1000"

	// Error that is caused by input
	StatusCodeGenericBadRequestError    = "2000"
	StatusCodeJsonNotParsableError      = "2001"
	StatusCodeMissingRequiredParameters = "2100"
	StatusCodeInvalidParameters         = "2200"
	StatusCodeDuplicatedEntry           = "2300"
	StatusCodeGenericNotFoundError      = "2500"
	StatusCodeUnprocessableEntity       = "2600"

	// Error that is caused by our own code
	StatusCodeGenericInternalError = "5000"
	StatusCodeDatabaseError        = "5100"

	// Error that is related to security
	StatusCodeAuthError    = "9000"
	StatusCodeUnauthorized = "9100"
	StatusCodeForbidden    = "9200"
	StatusCodeTokenExpired = "9300"
)
