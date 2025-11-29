package main
import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	cfg := &LocationConfig{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		cmd, ok := getCommands()[input]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(cfg); err != nil {
			fmt.Println("Error:", err)
		}
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


