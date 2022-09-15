package main

import "fmt"

func main() {

	// var a *int
	// *a += 1
	DoOperation(2, decrease)
}

func DoOperation(y int, f func(int, int)) {
	f(y, 1)
}

func increase(a, b int) int {

	return a + b
}

func decrease(a, b int) {
	fmt.Println("result is a- b ", a-b)
}
