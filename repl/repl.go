package repl

import (
	"bufio"
	"os"
	"fmt"
)



const prompt = "> "

func dispatch(command *string){

}

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

		dispatch(&text)
		fmt.Print(prompt)

	}
}
