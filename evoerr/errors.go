package evoerr

type InvalidLengthError struct{}

func (ile InvalidLengthError) Error() string {
	return "The length given was less than or equal to zero"
}
