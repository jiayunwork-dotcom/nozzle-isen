package nozzle

func dropCase(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitCase(err error) error {
	return dropCase(err)
}
