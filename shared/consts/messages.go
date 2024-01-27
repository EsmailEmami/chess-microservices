package consts

const (
	// Messages
	Created              = "New record created successfully"
	Updated              = "Update performed successfully"
	Deleted              = "Record deleted successfully"
	VerifivationCodeSent = "Verification code sent successfully."
	RegistrationDone     = "Registration completed successfully."
	LoggedOut            = "Logged out successfully."
	OperationDone        = "Operation completed successfully"

	// Errors
	CacheRecordNotFound = "Cache Record Not Found"
	BadRequest          = "Invalid data provided"
	InvalidRequest      = "Invalid request"
	InternalServerError = "An internal server error occurred"
	UnauthorizedError   = "Authentication error"
	ForbiddenError      = "Unauthorized access to the requested resource."
	ValidationError     = "Data validation error"
	Required            = "This field cannot be empty."
	InvalidValue        = "Invalid value entered."
	RecordNotFound      = "Requested record not found."
	PasswordIsShort     = "The password must be at least 8 characters long and include lowercase letters, uppercase letters, and special characters."
	InvalidCharacters   = "The entered value contains invalid characters."
)
