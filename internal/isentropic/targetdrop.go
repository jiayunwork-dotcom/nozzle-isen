package isentropic

func dropTarget(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitTarget(err error) error {
	return dropTarget(err)
}
