package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"ship-it/internal/api/models"

	"github.com/alecthomas/jsonschema"
)

func writeJSONSchema(name string, v interface{}) error {
	file, err := os.Create(fmt.Sprintf("api/%s.json", name))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(jsonschema.Reflect(v))
}

func main() {
	// NOTE: After creating a new JSON model type, it must be added to this
	// map in order to auto-generate JSON schema documents.
	types := map[string]interface{}{
		"release": models.Release{},
	}

	for name, typ := range types {
		if err := writeJSONSchema(name, typ); err != nil {
			log.Printf("failed to generate JSON schema doc for %s: %s\n", name, err.Error())
		}
	}
}
