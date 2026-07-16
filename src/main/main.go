package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	WorldHealthOrg "main/worldHealthOrg"
)

func printUsage() {
	fmt.Println("candybar CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  candybar help")
	fmt.Println("  candybar who-infant-nutrition <country>")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  candybar who-infant-nutrition USA")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	switch command {
	case "help", "-h", "--help":
		printUsage()
		return

	case "who-infant-nutrition":
		if len(os.Args) < 3 || strings.TrimSpace(os.Args[2]) == "" {
			fmt.Fprintln(os.Stderr, "country argument is required")
			printUsage()
			os.Exit(1)
		}

		country := strings.ToUpper(strings.TrimSpace(os.Args[2]))

		payload, err := WorldHealthOrg.FetchInfantNutrition(country)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch infant nutrition: %v\n", err)
			os.Exit(1)
		}

		encoded, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode output: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(encoded))
		return

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
