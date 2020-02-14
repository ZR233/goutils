package erroru

func warpedErr(err1 error) (err error) {

	err = AddInfo("wrap error1", err1)
	return
}
