package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/griffithsh/squads/game"
)

/*
char-var generates variations of character appearances
*/

type recolor struct {
	Name   string
	Colors []string
}
type config struct {
	HairColors []recolor
	SkinColors []recolor
}

type seedCloseup struct {
	X, Y int
}
type seed struct {
	Profession string
	Sexes      []string
	Input      string
	Closeup    seedCloseup
}

func main() {
	b, err := ioutil.ReadFile("./character-appearance/config.json")
	if err != nil {
		fmt.Printf("read config: %v", err)
		os.Exit(1)
	}
	var cfg config
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&cfg); err != nil {
		fmt.Printf("decode config: %v", err)
		os.Exit(1)
	}

	// read every file in ./character-appearance looking for .seed.json files
	files, err := ioutil.ReadDir("./character-appearance")
	if err != nil {
		fmt.Printf("read directory: %v", err)
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".seed.json") {
			// got a seed to use!
			b, err := ioutil.ReadFile("./character-appearance/" + file.Name())
			if err != nil {
				fmt.Printf("read seed %s: %v", file.Name(), err)
				os.Exit(1)
			}
			var s seed
			if err := json.NewDecoder(bytes.NewReader(b)).Decode(&s); err != nil {
				fmt.Printf("decode seed %s: %v", file.Name(), err)
				os.Exit(1)
			}

			imageFile := fmt.Sprintf("character-appearance/%s.variations%s", strings.TrimSuffix(s.Input, filepath.Ext(s.Input)), filepath.Ext(s.Input))
			inBytes, err := ioutil.ReadFile("character-appearance/" + s.Input)
			if err != nil {
				fmt.Printf("read file %s: %v", "character-appearance/"+s.Input, err)
				os.Exit(1)
			}

			// for every sex, haircolor and skincolor ...
			for _, sex := range s.Sexes {
				for fy, hair := range cfg.HairColors {
					for fx, skin := range cfg.SkinColors {
						// output an .appearance file
						appearanceFile := fmt.Sprintf("%s-%s(%s-%s).appearance", s.Profession, sex, hair.Name, skin.Name)

						var v struct {
							game.Appearance
							Profession string
							Sex        string
							HairColor  string
							SkinColor  string
						}
						v.Profession = s.Profession
						v.Sex = sex
						v.HairColor = hair.Name
						v.SkinColor = skin.Name
						v.Appearance.Participant = game.Sprite{
							Texture: imageFile,
							X:       0,
							Y:       0,
							W:       72,
							H:       96,
						}
						v.Appearance.Portrait = game.Sprite{
							Texture: imageFile,
							X:       s.Closeup.X - 13,
							Y:       s.Closeup.Y - 13,
							W:       26,
							H:       26,
						}
						buf := bytes.Buffer{}
						if err := json.NewEncoder(&buf).Encode(&v); err != nil {
							fmt.Printf("encode appearance %s: %v", appearanceFile, err)
							os.Exit(1)
						}
						ioutil.WriteFile("character-appearance/"+appearanceFile, buf.Bytes(), 0644)

						// TODO: replace pixels in a copy of the input image and
						// composite it into a patchwork.
						fmt.Printf("(%s-%s) this variation needs to go at %d,%d\n", hair.Name, skin.Name, fx, fy)
					}
				}
			}
			// write out the file here ...
			if err := ioutil.WriteFile(imageFile, inBytes, 0644); err != nil {
				fmt.Printf("write file %s: %v", imageFile, err)
				os.Exit(1)
			}
		}
	}
}
