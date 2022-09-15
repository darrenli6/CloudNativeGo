package main

import (
	"fmt"
)

/*
interface 是一种事物的抽象


*/

type IF interface {
	getName() string
}

type Human struct {
	firstName, lastName string
}

type Plant struct {
	vendor string
	model  string
}

func (h *Human) getName() string {
	return h.firstName + "," + h.lastName
}

func (p Plant) getName() string {
	return fmt.Sprintf("vendor : %s  model : %s", p.vendor, p.model)
}

type Car struct {
	factory, model string
}

func (c *Car) getName() string {
	return c.factory + "-" + c.model
}

func main() {
	interfaces := []IF{}
	h := new(Human)
	h.firstName = "darren"
	h.lastName = "li"
	interfaces = append(interfaces, h)

	s := new(Car)
	s.factory = "byd"
	s.model = "dmi"

	interfaces = append(interfaces, s)

	for _, f := range interfaces {
		fmt.Println(f.getName())
	}

	p := Plant{}
	p.model = "p"
	p.vendor = "pv"
	fmt.Println(p)

}
