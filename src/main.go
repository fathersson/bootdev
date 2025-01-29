package main

import (
	"fmt"
	"strings"
)

func main() {
	//cleanInput("  hello  world  ")
	//fmt.Println(result)
	//fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	var c int
	var result []string
	for i := 0; i < len(text); i++ {
		if text[i] != ' ' {
			var str string
			for _, s := range text[i:] {
				if s == ' ' {
					break
				}
				if strings.ToUpper(string(s)) == string(s) {
					symbol := strings.ToLower(string(s))
					str = str + symbol
					c++
					continue
				}
				c++
				str = str + string(s)
			}
			result = append(result, str)
			str = ""
			i = i + c
			c = 0

		}

	}
	fmt.Println(result)
	return result
}
