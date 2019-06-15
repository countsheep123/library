package main

func before() error {
	if err := validate(); err != nil {
		return err
	}

	return nil
}
