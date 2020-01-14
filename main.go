package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(os.Args) != 4 {
		fmt.Println("Usage:", os.Args[0], "<site>", "<budget>", "<bedrooms>")
		return
	}
	site := args[1]
	budget := args[2]
	bedrooms := args[3]

	a := &App{}
	err := a.Init("config.yml", site)
	if err != nil {
		log.Fatal(err)
	}
	a.Watch(budget, bedrooms)
}
