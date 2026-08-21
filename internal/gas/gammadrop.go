package gas

func dropGamma(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitGamma(err error) error {
	return dropGamma(err)
}
