package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("rlp", "Read relationship properties")

	app.Spec = "NAME [PROPERTY]"

	var (
		relName  = app.StringArg("NAME", "", "Relationship name")
		property = app.StringArg("PROPERTY", "", "Relationship property")
	)

	app.Action = func() {
		var relationships map[string][]map[string]interface{}

		relationshipsString := os.Getenv("PLATFORM_RELATIONSHIPS")
		if relationshipsString == "" {
			fmt.Println("No relationships found")
			os.Exit(1)
		}

		relationshipsBytes, _ := base64.StdEncoding.DecodeString(relationshipsString)
		if err := json.Unmarshal(relationshipsBytes, &relationships); err != nil {
			log.Fatal(err)
		}

		if list, exists := relationships[*relName]; exists {
			for _, relationship := range list {
				if *property == "" {
					formatted, _ := json.MarshalIndent(relationship, "", "  ")
					fmt.Printf("%s", formatted)
				} else {
					fmt.Printf("%s", relationship[*property])
				}
			}
		} else {
			fmt.Printf("Relationship not found: %s\n", *relName)
			os.Exit(1)
		}

	}

	app.Run(os.Args)
}
