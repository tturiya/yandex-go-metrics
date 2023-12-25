package main

import "github.com/tturiya/iter5/internal/server"

func main() {
	err := server.StartServer()
	if err != nil {
		panic(err)
	}
}
