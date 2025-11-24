package main
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Print("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)

	map[string] cliCommand {
    	"exit": {
        	name:        "exit",
      	  description: "Exit the Pokedex",
    	    callback:    commandExit,
  	  },
	}

	for scanner.Scan() {
		if cliCommand[scanner.Text()] != nil{
			return cliCommand[scanner.Text()].callback
		}
		fmt.Print("Pokedex > ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}


func cleanInput(text string) []string {
  slice := []string{}
  if text == ""{
		return slice
	}
	slice = strings.Fields(strings.ToLower(text))
	return slice
}


type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


