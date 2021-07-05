package types

type KongError struct {
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}
