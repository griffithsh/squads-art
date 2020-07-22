package main

import (
	"flag"
	"fmt"
)

// Accepts a directory of exported tmx and tsx files as json, and converts them
// to the squads game.CombatMapRecipe type serialised as json.
// The output is a folder of PNGs and Marshaled game.CombatMapRecipes

func main() {

	inDir := flag.String("in", "./", "configures the directory to scan through looking for *.tmx.json and *.tsx.json")
	outDir := flag.String("out", "./", "configures the output directory to write *.terrain and dependant *.png files to")
	flag.Parse()

	c := newConverter()
	err := c.scan(*inDir)
	if err != nil {
		panic(fmt.Sprintf("scan: %v", err))
	}

	err = c.convert()
	if err != nil {
		panic(fmt.Sprintf("convert: %v", err))
	}

	err = c.write(*outDir)
	if err != nil {
		panic(fmt.Sprintf("write: %v", err))
	}
}
