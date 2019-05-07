package common

//errorWithCodeData error dengan error code + message
type errorWithCodeData struct {
	//ErrorMessage error message. human readable
	ErrorMessage string
	//ErrorCode business / system code for error
	ErrorCode string
	//RawError original raw error data
	RawError error
}

//ErrorWithCodeData error with code
type ErrorWithCodeData interface {
	error
	//GetErrorCode read error code
	GetErrorCode() string
	//GetRawError rawerror data
	GetRawError() error
}

//ErrorWithCode generate error with code
func ErrorWithCode(errorMessage string, errorCode string) ErrorWithCodeData {
	return &errorWithCodeData{ErrorMessage: errorMessage, ErrorCode: errorCode}
}

//ErrorWithCodeAndRawError generate error with raw error
func ErrorWithCodeAndRawError(errorMessage string, errorCode string, err error) ErrorWithCodeData {
	return &errorWithCodeData{ErrorMessage: errorMessage, ErrorCode: errorCode, RawError: err}
}

//Error display error string( hanya message saja di render. ikut mekanisme standard)
func (errActual *errorWithCodeData) Error() string {
	return errActual.ErrorMessage
}

//GetErrorCode get error code of data
func (errActual *errorWithCodeData) GetErrorCode() string {
	return errActual.ErrorCode
}
func (errActual *errorWithCodeData) GetRawError() error {
	return errActual.RawError
}
