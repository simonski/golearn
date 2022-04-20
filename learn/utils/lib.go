package utils

import "fmt"

func Foo() string {
	return "Foo"
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("error!")
		panic(err)
	}
}
