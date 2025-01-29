package main

import (
	"fmt"
	"pattern/pattern"
)

func main() {
	fmt.Println("=== Testing Facade Pattern ===")
	cinema := pattern.NewCinema()
	cinema.StartCinema()

}
