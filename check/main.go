package main

import (
	"fmt"
)

func main() {
	output, err := Execute()
	if err != nil {
		panic(err)
	}
	fmt.Println(output)
}

//Execute - provides check capability
func Execute() (string, error) {
	return "[]", nil
}
