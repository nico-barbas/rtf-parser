package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	input, err := os.ReadFile("./test.rtf")
	if err != nil {
		log.Fatal(err)
		return
	}

	ops, err := Parse(string(input))

	if err != nil {
		fmt.Println(err)
	} else {
		PrettyPrintEntities(ops)

		output := BuildLayoutHTML(ops)
		fmt.Println(output)
	}
}
