package repl

import (
	"bufio"
	"os"
	"fmt"
)

const prompt = "> "

func StartInteractive(){

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)

	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
		fmt.Println(len(text))

		if "q" == text|| "quit" == text{
			break
		}

		fmt.Print(prompt)

	}
}
