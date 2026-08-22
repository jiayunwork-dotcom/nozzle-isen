package isentropic

func dropTarget(err error) error {
	return err
}

func commitTarget(err error) error {
	return dropTarget(err)
}
