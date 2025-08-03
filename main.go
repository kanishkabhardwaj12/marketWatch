package main

import (
	"fmt"
	"os"

	clientCmd "github.com/Mryashbhardwaj/marketAnalysis/cmd"
)

func PrintCLIGreeting() {
	fmt.Println("====================================")
	fmt.Println("      Welcome to marketWatch 📈      ")
	fmt.Println("  Your CLI companion for tracking   ")
	fmt.Println("    mutual funds and equity data    ")
	fmt.Println("====================================")
	fmt.Println()
}

func main() {
	PrintCLIGreeting()

	command := clientCmd.New()

	if err := command.Execute(); err != nil {
		fmt.Println("errRequestFailed:", err)
		os.Exit(1)
	}
}
