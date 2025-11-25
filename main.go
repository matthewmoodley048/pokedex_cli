package main
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	fmt.Print("Pokedex > ")
	scanner := bufio.NewScanner(os.Stdin)

 commands:=	map[string] cliCommand {
    	"exit": {
        	name:        "exit",
      	  description: "Exit the Pokedex",
    	    callback:    commandExit,
  	  },
			"help": {
					name:				 "help",
					description: "Displays a help message",
					callback: 	 commanHelp,
			},
	}

	for scanner.Scan() {
		if c,ok:= commands[scanner.Text()]; ok{
	    if err := c.callback(); err != nil {
      	  fmt.Println(err)
  	  }
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


func commanHelp() error {
	fmt.Println("Welcome to the Pokedex! \nUsage:\n")	
	commands:=	map[string] cliCommand {
    	"exit": {
        	name:        "exit",
      	  description: "Exit the Pokedex",
    	    callback:    commandExit,
  	  },
			"help": {
					name:				 "help",
					description: "Displays a help message",
					callback: 	 commanHelp,
			},
	}

	for _, command := range commands{
		fmt.Printf("%s: %s \n", command.name, command.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}


