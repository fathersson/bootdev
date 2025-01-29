package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//cleanInput("  hello  world  ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		firstWord := len(words[0]) + 1

		var str string
		var fWord string
		for i, word := range words {
			for _, s := range word {
				if strings.ToUpper(string(s)) == string(s) {
					symbol := strings.ToLower(string(s))
					str = str + symbol

					firstWord--
					if firstWord > 0 {
						fWord = fWord + symbol
					}
					continue
				}
				str = str + string(s)
				firstWord--
				if firstWord > 0 {
					fWord = fWord + string(s)
				}
			}
			if i < len(words)-1 {
				str = str + " "
			}
		}
		//fmt.Print("Pokedex > ", str)
		//fmt.Print("\n")
		fmt.Print("Your command was: ", fWord)
		fmt.Print("\n")
		//fmt.Print(fWord)
	}
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
