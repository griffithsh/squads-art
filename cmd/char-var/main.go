package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
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
	Feet       seedCloseup
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
			input, _, err := image.Decode(bytes.NewReader(inBytes))
			bounds := input.Bounds()
			w := (bounds.Max.X - bounds.Min.X)
			h := (bounds.Max.Y - bounds.Min.Y)
			canvas := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w * len(cfg.SkinColors), h * len(cfg.HairColors)}})

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
							X:       fx * w,
							Y:       fy * h,
							W:       w,
							H:       h,
						}
						v.Appearance.FaceX = s.Closeup.X
						v.Appearance.FaceY = s.Closeup.Y

						v.Appearance.FeetX = s.Feet.X
						v.Appearance.FeetY = s.Feet.Y

						buf := bytes.Buffer{}
						if err := json.NewEncoder(&buf).Encode(&v); err != nil {
							fmt.Printf("encode appearance %s: %v", appearanceFile, err)
							os.Exit(1)
						}
						ioutil.WriteFile("character-appearance/"+appearanceFile, buf.Bytes(), 0644)

						input, _, err := image.Decode(bytes.NewReader(inBytes))
						if err != nil {
							panic(fmt.Sprintf("decode input image again: %v", err))
						}

						// TODO: replace pixels in a copy of the input image.
						// ...

						// Composite the copy into the canvas of variations.
						destRect := image.Rect(fx*w, fy*h, fx*w+w, fy*h+h)
						draw.Draw(canvas, destRect, input, image.Point{0, 0}, draw.Over)
					}
				}
			}
			// write out the file here ...
			f, err := os.Create(imageFile)
			if err != nil {
				panic(fmt.Sprintf("create canvas file %s: %v", imageFile, err))
			}
			err = png.Encode(f, canvas)
			if err != nil {
				panic(fmt.Sprintf("encode png: %v", err))
			}
		}
	}
}
