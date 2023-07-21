package main

import (
	"fmt"

	"github.com/Studiumz/studiumz-api/app"
)

func main() {
	c := app.LoadConfig()
	fmt.Println(c)
}
