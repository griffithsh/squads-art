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
	"path"
	"path/filepath"
	"strings"

	"github.com/griffithsh/squads/game"
)

/*
char-var generates variations of character appearances
*/

type point struct {
	X, Y int
}
type seed struct {
	Profession string
	Sexes      []string
	Input      string
	Closeup    point
	Feet       point
}

func main() {
	b, err := ioutil.ReadFile("./character-appearance/config.json")
	if err != nil {
		fmt.Printf("read config: %v\n", err)
		os.Exit(1)
	}
	var cfg config
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&cfg); err != nil {
		fmt.Printf("decode config: %v\n", err)
		os.Exit(1)
	}

	// read every file in ./character-appearance looking for .seed.json files
	files, err := ioutil.ReadDir("./character-appearance")
	if err != nil {
		fmt.Printf("read directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".seed.json") {
			// got a seed to use!
			b, err := ioutil.ReadFile("./character-appearance/" + file.Name())
			if err != nil {
				fmt.Printf("read seed %s: %v\n", file.Name(), err)
				os.Exit(1)
			}
			var s seed
			if err := json.NewDecoder(bytes.NewReader(b)).Decode(&s); err != nil {
				fmt.Printf("decode seed %s: %v\n", file.Name(), err)
				os.Exit(1)
			}

			imageFile := fmt.Sprintf("character-appearance/%s.variations%s", strings.TrimSuffix(s.Input, filepath.Ext(s.Input)), filepath.Ext(s.Input))
			inBytes, err := ioutil.ReadFile("character-appearance/" + s.Input)
			if err != nil {
				fmt.Printf("read file %s: %v\n", "character-appearance/"+s.Input, err)
				os.Exit(1)
			}
			input, _, err := image.Decode(bytes.NewReader(inBytes))
			if err != nil {
				fmt.Printf("decode bytes %s: %v\n", "character-appearance/"+s.Input, err)
				os.Exit(1)
			}
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
						if err := ioutil.WriteFile("packed/content/character-appearance/"+appearanceFile, buf.Bytes(), 0644); err != nil {
							panic(fmt.Sprintf("writing file %v", err))
						}

						raw, _, err := image.Decode(bytes.NewReader(inBytes))
						input := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
						draw.Draw(input, image.Rect(0, 0, w, h), raw, image.Point{0, 0}, draw.Over)

						if err != nil {
							panic(fmt.Sprintf("decode input image again: %v", err))
						}

						rm := remapper{}
						rm.load(cfg.HairTemplate, hair.Colors)
						rm.load(cfg.SkinTemplate, skin.Colors)

						// For every pixel of the image ...
						for i := 0; i < w*h; i++ {
							x, y := i%w, i/w
							curr := input.At(x, y)
							_, _, _, a := curr.RGBA()
							if a == 0 {
								continue
							}
							if should, to := rm.remap(curr); should {
								input.Set(x, y, to)
							}
						}

						// Composite the copy into the canvas of variations.
						destRect := image.Rect(fx*w, fy*h, fx*w+w, fy*h+h)
						draw.Draw(canvas, destRect, input, image.Point{0, 0}, draw.Over)
					}
				}
			}

			// write out the file here ...
			f, err := os.Create(path.Join("packed", "content", imageFile))
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
