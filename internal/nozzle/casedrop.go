package nozzle

func dropCase(err error) error {
	return err
}

func commitCase(err error) error {
	return dropCase(err)
}
