package errort

func warpedErr(err1 error) (err error) {

	err = NewFromError(err1, 0)
	return
}
