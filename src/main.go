package main

import (
	"fmt"

	"github.com/jumaevkova04/search/pkg/search"
)

func main() {
	// 	a := search.FindAllPhraseInFile("name", "src/text.txt")
	// 	fmt.Println(a)
	b := search.FindAnyPhraseInFile("name", "src/text.txt")
	fmt.Println(b)
}
