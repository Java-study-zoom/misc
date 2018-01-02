package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func main() {
	envs := os.Environ()
	sort.Strings(envs)
	bs, err := json.MarshalIndent(envs, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(bs); err != nil {
		log.Fatal(err)
	}
}
