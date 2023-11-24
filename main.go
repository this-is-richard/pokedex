package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/this-is-richard/pokedex/pokedexcli"
)

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	p := pokedexcli.NewPokedex()

	fmt.Printf("pokedex > ")

	for scanner.Scan() {
		text := scanner.Text()
		words := strings.Split(text, " ")
		commandKey, rest := words[0], words[1:]

		err := p.Run(commandKey, rest...)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("pokedex > ")
	}
}

func main() {
	Repl()
}
