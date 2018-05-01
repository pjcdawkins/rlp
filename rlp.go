package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jawher/mow.cli"
)

func errorOut(message string) {
	os.Stderr.WriteString(message + "\n")
	os.Exit(1)
}

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
			errorOut("The PLATFORM_RELATIONSHIPS variable is not defined")
		}

		relationshipsBytes, _ := base64.StdEncoding.DecodeString(relationshipsString)
		if err := json.Unmarshal(relationshipsBytes, &relationships); err != nil {
			errorOut(err.Error())
		}

		if len(relationships) == 0 {
			errorOut("No relationships found")
		}

		relationshipsList, relationshipExists := relationships[*relName]
		if !relationshipExists {
			os.Stderr.WriteString(fmt.Sprintf("Relationship not found: '%s'\n", *relName))
			if len(relationships) > 0 {
				os.Stderr.WriteString("Available relationships:\n")
				for name, _ := range relationships {
					os.Stderr.WriteString("  " + name + "\n")
				}
			}
			os.Exit(1)
		}

		for _, relationship := range relationshipsList {
			if *property == "" {
				formatted, _ := json.MarshalIndent(relationship, "", "  ")
				fmt.Printf("%s\n", formatted)
				continue
			}

			if value, propertyExists := relationship[*property]; propertyExists {
				fmt.Printf("%v\n", value)
				continue
			}

			errorOut(fmt.Sprintf("Property '%s' not found in relationship '%s'\n", *property, *relName))
		}
	}

	app.Run(os.Args)
}
