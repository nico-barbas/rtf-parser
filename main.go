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

		layout := BuildLayout(ops)
		fmt.Printf("%#v\n", layout)

		output := OutputHTML(layout, BuilderOptions{prettyOutput: false})
		fmt.Println(output)

		outputFile, err := os.OpenFile("./output.html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		_, err = outputFile.WriteString(output)
		if err != nil {
			log.Fatal(err)
		}
	}
}
