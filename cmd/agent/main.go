package main

import "github.com/tturiya/iter5/internal/agent"

func main() {
	err := agent.StartAgent()
	if err != nil {
		panic(err)
	}
}
