package common

func RaiseWhenNok(ok bool, err error) error {
	if ok {
		return nil
	}
	return err
}
