package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"

	"github.com/nei7/ntube/internal/rest"
	"gopkg.in/yaml.v2"
)

func main() {
	var output string
	flag.StringVar(&output, "path", "", "")
	flag.Parse()

	if output == "" {
		log.Fatal("output is required")
	}

	swagger := rest.NewOpenAPI3()

	data, err := json.Marshal(&swagger)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(path.Join(output, "openapi3.json"), data, 0600); err != nil {
		log.Fatal(err)
	}

	data, err = yaml.Marshal(&swagger)
	if err != nil {
		log.Fatalf("Couldn't marshal json: %s", err)
	}

	if err := os.WriteFile(path.Join(output, "openapi3.yaml"), data, 0600); err != nil {
		log.Fatalf("Couldn't write json: %s", err)
	}

}
