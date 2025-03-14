package response

const (
	ErrCodeSucsess = 2001
	ErrCodeParamInvalid = 20003 // email invalid
)

var msg = map[int]string {
	ErrCodeSucsess: "Success",
	ErrCodeParamInvalid: "Email is invalid",
}
