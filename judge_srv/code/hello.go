package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	a := ""
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Print(a)
			break
		}
		a += text
	}
}
