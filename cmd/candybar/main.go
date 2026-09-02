package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"candybar/internal/food"
	"candybar/internal/media"
	"candybar/internal/worldhealthorg"
)

func printUsage() {
	fmt.Println("candybar CLI")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  candybar help")
	fmt.Println("  candybar who-infant-nutrition <country>")
	fmt.Println("  candybar who-infant-deaths <country>")
	fmt.Println("  candybar food-recalls <product description>")
	fmt.Println("  candybar compare-images [file]")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  candybar who-infant-nutrition USA")
	fmt.Println("  candybar who-infant-deaths USA")
	fmt.Println("  candybar food-recalls peanut")
	fmt.Println("  candybar compare-images figures.json")
	fmt.Println("  cat figures.json | candybar compare-images")
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

		payload, err := worldhealthorg.FetchInfantNutrition(country)
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

	case "who-infant-deaths":
		if len(os.Args) < 3 || strings.TrimSpace(os.Args[2]) == "" {
			fmt.Fprintln(os.Stderr, "country argument is required")
			printUsage()
			os.Exit(1)
		}

		country := strings.ToUpper(strings.TrimSpace(os.Args[2]))

		payload, err := worldhealthorg.FetchInfantDeaths(country)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch infant deaths: %v\n", err)
			os.Exit(1)
		}

		encoded, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode output: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(encoded))
		return

	case "food-recalls":
		if len(os.Args) < 3 || strings.TrimSpace(strings.Join(os.Args[2:], " ")) == "" {
			fmt.Fprintln(os.Stderr, "product description argument is required")
			printUsage()
			os.Exit(1)
		}

		productDescription := strings.TrimSpace(strings.Join(os.Args[2:], " "))

		payload, err := food.FetchRecallsByProductDescription(productDescription)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to fetch food recalls: %v\n", err)
			os.Exit(1)
		}

		encoded, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode output: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(encoded))
		return

	case "compare-images":
		var input io.Reader = os.Stdin

		if len(os.Args) >= 3 && strings.TrimSpace(os.Args[2]) != "" {
			file, err := os.Open(os.Args[2])
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()

			input = file
		}

		var figures []media.Figure
		if err := json.NewDecoder(input).Decode(&figures); err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode input: %v\n", err)
			os.Exit(1)
		}

		result := media.Compare(figures)

		encoded, err := json.MarshalIndent(result, "", "  ")
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
