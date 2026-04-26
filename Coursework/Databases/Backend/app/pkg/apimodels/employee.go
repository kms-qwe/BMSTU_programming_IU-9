package apimodels

type Employee struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Login    string `json:"login"`
	Phone    string `json:"phone"`
}

type ListEmployeesResponse struct {
	Employees []Employee `json:"employees"`
}
