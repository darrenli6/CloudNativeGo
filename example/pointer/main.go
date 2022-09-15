package main

import "fmt"

/*
面试

go语言 到底传值还是传指针

**/

func main() {
	str := "i love you "
	fmt.Println(str)
	point := &str

	fmt.Println(point)

	anotherString := *&str
	fmt.Println(anotherString)

	str = "i love me"
	fmt.Println(str)
	fmt.Println(point)
	fmt.Println(anotherString)

	param := ParameterStruct{Name: "aaa"}

	fmt.Println(param)

	changeParameter(&param, "bbbb")
	fmt.Println(param)

	cantChangeParameter(param, "ddd")
	fmt.Println(param)

}

type ParameterStruct struct {
	Name string
}

func changeParameter(param *ParameterStruct, value string) {
	param.Name = value
}

func cantChangeParameter(param ParameterStruct, value string) {
	param.Name = value
}
