package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{Id: "1001", Name: "John Doe", City: "New York", Zipcode: "10001", Birthdate: "1985-01-15", Status: "1"},
		{Id: "1002", Name: "Jane Smith", City: "Los Angeles", Zipcode: "90001", Birthdate: "1990-05-20", Status: "1"},
		{Id: "1003", Name: "Peter Jones", City: "Chicago", Zipcode: "60601", Birthdate: "1982-11-30", Status: "0"},
		{Id: "1004", Name: "Susan Williams", City: "Houston", Zipcode: "77001", Birthdate: "1995-08-02", Status: "1"},
		{Id: "1005", Name: "David Brown", City: "Phoenix", Zipcode: "85001", Birthdate: "1988-03-25", Status: "1"},
	}
	return CustomerRepositoryStub{customers: customers}
}
