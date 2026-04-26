package models

type Employee struct {
	ID       string
	FullName string
	Login    string
	Phone    string
}

type EmployeeCredentials struct {
	Employee
	PasswordHash string
}
