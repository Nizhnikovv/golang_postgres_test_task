package main

type CreateEmployeeRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
}

func (r *CreateEmployeeRequest) ToEmployee() *Employee {
	return &Employee{
		FirstName:  r.FirstName,
		LastName:   r.LastName,
		MiddleName: r.MiddleName,
		Phone:      r.Phone,
		City:       r.City,
	}
}

type Employee struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
}
