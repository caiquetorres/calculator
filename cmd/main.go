package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/caiquetorres/calculator/eval"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for reader.Scan() {
		input := bytes.NewReader(reader.Bytes())
		res, err := eval.Eval(input)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(fmt.Sprintf("%f", res))
		}
		fmt.Print("> ")
	}
}
