package gevaerr

type InvalidLengthError struct{}

func (ile InvalidLengthError) Error() string {
	return "The length given was less than or equal to zero"
}

type InvalidParamError struct {
	ParamName string
}

func (ipe InvalidParamError) Error() string {
	return "The parameter " + ipe.ParamName + " was invalid"
}
