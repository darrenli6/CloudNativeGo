package main

import "fmt"

func changeSlice(myArr []string) {

	for k, v := range myArr {
		if v == "stupid" {
			myArr[k] = "smart"
		} else if v == "weak" {
			myArr[k] = "strong"
		}

	}

}

func main() {

	myArr := []string{"I", "am", "stupid", "and", "weak"}

	fmt.Printf("myAarr is %v \n ", myArr)

	changeSlice(myArr)

	fmt.Printf("myAarr is %v \n ", myArr)
	for _, v := range myArr {
		fmt.Println("", v)
	}
}
