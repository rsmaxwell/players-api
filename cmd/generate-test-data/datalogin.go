package main

// createTestdataLoggedin function
func createTestdataLoggedin() error {

	err := clearModel()
	if err != nil {
		return err
	}

	datapeople := []struct {
		id        string
		password  string
		firstname string
		lastname  string
		email     string
	}{
		{id: "007", password: "topsecret", firstname: "James", lastname: "Bond", email: "james@mi6.co.uk"},
	}

	for _, i := range datapeople {
		err = Register(i.id, i.password, i.firstname, i.lastname, i.email)
		if err != nil {
			return err
		}
	}

	return nil
}
