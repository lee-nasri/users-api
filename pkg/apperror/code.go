package apperror

type (
	ErrorCode int
)

const (
	Success ErrorCode = 0

	ErrInvalidRequest ErrorCode = 6000

	ErrUserIDAleadyExist ErrorCode = 7001
	ErrUserNotFound      ErrorCode = 7010

	ErrInternalServerError ErrorCode = 9000
)

var (
	Message = map[ErrorCode]string{
		// 600x
		ErrInvalidRequest: "Invalid request",

		// 700x
		ErrUserIDAleadyExist: "User ID already exist in database",
		ErrUserNotFound:      "User Not Found",

		// 900x
		ErrInternalServerError: "Internal Server Error",
	}
)
