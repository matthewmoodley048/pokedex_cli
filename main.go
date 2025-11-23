package main
import (
	"fmt";
  "strings"
)

func main() {
}

func cleanInput(text string) []string {
  slice := []string{}
  if text == ""{
		return slice
	}
	slice = strings.Fields(strings.ToLower(text))
	return slice
}
