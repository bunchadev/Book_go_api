package response

const (
	ErrCodeSuccess      = 300
	ErrCodeParamInvalid = 301
	ErrCodeServer       = 302
	ErrCodeUserName     = 303
	ErrCodeLogin        = 304
	ErrCodeAuth         = 305
)

var msg = map[int]string{
	ErrCodeSuccess:      "success",
	ErrCodeParamInvalid: "Param invalid",
	ErrCodeServer:       "Faulty server",
	ErrCodeUserName:     "Existing accounts",
	ErrCodeLogin:        "Username or password invalid",
	ErrCodeAuth:         "Not Authorized",
}
