package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeId int
}

func (e *Employee) PrintInfo() {
	fmt.Println("Name:", e.Name, "\tAge:", e.Age, "\tEmployeeId:", e.EmployeeId)
}

func main() {
	e := new(Employee)
	e.Name = "Jack"
	e.Age = 25
	e.EmployeeId = 1

	e.PrintInfo()
}
