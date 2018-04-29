package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
			os.Stderr.WriteString("No relationships found\n")
			os.Exit(1)
		}

		relationshipsBytes, _ := base64.StdEncoding.DecodeString(relationshipsString)
		if err := json.Unmarshal(relationshipsBytes, &relationships); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		relationshipsList, relationshipExists := relationships[*relName]
		if !relationshipExists {
			os.Stderr.WriteString(fmt.Sprintf("Relationship not found: '%s'\n", *relName))
			os.Exit(1)
		}

		for _, relationship := range relationshipsList {
			if *property == "" {
				formatted, _ := json.MarshalIndent(relationship, "", "  ")
				fmt.Printf("%s", formatted)
				continue
			}

			if value, propertyExists := relationship[*property]; propertyExists {
				fmt.Printf("%v", value)
				continue
			}

			os.Stderr.WriteString(fmt.Sprintf("Property '%s' not found in relationship '%s'\n", *property, *relName))
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
