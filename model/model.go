package model

// Person object
type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Address  *Address
}

// Address object
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}
