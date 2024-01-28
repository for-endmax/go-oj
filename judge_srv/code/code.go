package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	a := ""
	time.Sleep(time.Second * 10)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Print(a)
			break
		}
		a += text
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "读取标准输入时发生错误:", err)
	}
}
