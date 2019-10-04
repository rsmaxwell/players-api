package main

import "github.com/rsmaxwell/players-api/internal/model"

// createTestdataEmpty function
func createTestdataEmpty() error {

	err := model.ClearModel()
	if err != nil {
		return err
	}

	return nil
}
