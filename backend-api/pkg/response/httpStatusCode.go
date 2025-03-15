package response

const (
	ErrCodeSucsess = 2001
	ErrCodeParamInvalid = 20003 // email invalid
	TokenInvalid = 20004
)

var msg = map[int]string {
	ErrCodeSucsess: "Success",
	ErrCodeParamInvalid: "Email is invalid",
	TokenInvalid: "Token is invalid",
}
