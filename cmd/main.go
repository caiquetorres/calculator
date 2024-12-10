package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/caiquetorres/calculator/eval"
)

func main() {
	input := strings.NewReader("1 + 2")
	res, err := eval.Eval(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("Res: %f", res))
}
