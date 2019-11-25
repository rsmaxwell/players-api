package main

// createTestdataEmpty function
func createTestdataEmpty() error {

	err := clearModel()
	if err != nil {
		return err
	}

	return nil
}
