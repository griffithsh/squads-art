// aseprite-deseed is some pretty quickly hacked together code to merge some
// critical meta information about a profession's artwork, (like how fast it
// moves and where it's attack apex frame is) with some actual artwork and frame
// data exported from aseprite. See aseprite-export.sh in the root.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/griffithsh/squads/game"
)

type seed struct {
	From    string `json:"from"`
	OffsetX int    `json:"offsetX"`
	OffsetY int    `json:"offsetY"`
	game.PerformanceSet
}

func (s *seed) UnmarshalJSON(data []byte) error {
	var v struct {
		From    string
		OffsetX int
		OffsetY int
	}
	err := json.Unmarshal(data, &s.PerformanceSet)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	s.From = v.From
	s.OffsetX = v.OffsetX
	s.OffsetY = v.OffsetY
	return nil
}

type Data struct {
	X int
	Y int
	W int
	H int
}
type Frame struct {
	Filename string
	Data     Data `json:"frame"`
	Duration int
}
type Meta struct {
	Image string
}
type aseprite struct {
	Frames []Frame
	Meta   Meta
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	seeds := []seed{}

	dec.Decode(&seeds)

	for _, seed := range seeds {
		b, err := ioutil.ReadFile(seed.From)
		if err != nil {
			fmt.Println("ReadFile:", err)
			os.Exit(1)
		}
		var a aseprite
		json.Unmarshal(b, &a)

		// iterate through everything in a.Frames, patching it in to seed.PerformanceSet
		for _, frame := range a.Frames {
			gf := game.Frame{
				DurationMs: frame.Duration,
				Sprite: game.Sprite{
					Texture: a.Meta.Image,
					X:       frame.Data.X,
					Y:       frame.Data.Y,
					W:       frame.Data.W,
					H:       frame.Data.H,
					OffsetX: seed.OffsetX,
					OffsetY: seed.OffsetY,
				},
			}
			parts := strings.Split(frame.Filename, "-")
			animation := parts[0]
			dir := parts[1]
			num, _ := strconv.Atoi(parts[2])

			resolve := func(anim, dir string) *[]game.Frame {
				resolveDir := func(dir string, pfd *game.PerformancesForDirection) *[]game.Frame {
					switch dir {
					case "s":
						return &pfd.S
					case "n":
						return &pfd.N
					case "sw":
						return &pfd.SW
					case "se":
						return &pfd.SE
					case "nw":
						return &pfd.NW
					case "ne":
						return &pfd.NE
					}
					panic(fmt.Sprintf("resolveDir for invalid dir(%s)", dir))
				}
				switch anim {
				case "idle":
					return resolveDir(dir, &seed.PerformanceSet.Idle)
				case "move":
					return resolveDir(dir, &seed.PerformanceSet.Move)
				case "attack":
					return resolveDir(dir, &seed.PerformanceSet.Attack)
				case "spell":
					if dir != "s" {
						break
					}
					return &seed.PerformanceSet.Spell
				case "death":
					if dir != "s" {
						break
					}
					return &seed.PerformanceSet.Death
				case "rise":
					if dir != "s" {
						break
					}
					return &seed.PerformanceSet.Rise
				case "victory":
					if dir != "s" {
						break
					}
					return &seed.PerformanceSet.Victory

				}
				panic(fmt.Sprintf("There is no animation \"%s\" for direction \"%s\"", anim, dir))
			}

			f := resolve(animation, dir)
			if len(*f) != num {
				// This means that either there are out-of-order frames, or duplicate frames.
				panic(fmt.Sprintf("incorrect length (want %d, got %d) indicates data problem: animation(%s), direction(%s)", num, len(*f), animation, dir))
			}
			*f = append(*f, gf)
			continue
		}

		// go through seed.PerformanceSet, patching in a default animation for anything that has no frames
		defaultFrame := []game.Frame{game.Frame{
			DurationMs: 1000,
			Sprite: game.Sprite{
				Texture: a.Meta.Image,
				X:       0,
				Y:       0,
				W:       0,
				H:       0,
			},
		}}
		ensure := func(pfd *game.PerformancesForDirection) {
			if pfd.S == nil {
				pfd.S = defaultFrame
			}
			if pfd.N == nil {
				pfd.N = defaultFrame
			}
			if pfd.SE == nil {
				pfd.SE = defaultFrame
			}
			if pfd.SW == nil {
				pfd.SW = defaultFrame
			}
			if pfd.NE == nil {
				pfd.NE = defaultFrame
			}
			if pfd.NW == nil {
				pfd.NW = defaultFrame
			}
		}
		ensure(&seed.PerformanceSet.Idle)
		ensure(&seed.PerformanceSet.Move)
		ensure(&seed.PerformanceSet.Attack)
		if seed.PerformanceSet.Spell == nil {
			seed.PerformanceSet.Spell = defaultFrame
		}
		if seed.PerformanceSet.Death == nil {
			seed.PerformanceSet.Death = defaultFrame
		}
		if seed.PerformanceSet.Rise == nil {
			seed.PerformanceSet.Rise = defaultFrame
		}
		if seed.PerformanceSet.Victory == nil {
			seed.PerformanceSet.Victory = defaultFrame
		}

		// Write out seed.PerformanceSet.
		b, err = json.Marshal(seed.PerformanceSet)
		ioutil.WriteFile(fmt.Sprintf("%s.performance-set", seed.Name), b, 0644)

		err = os.Remove(seed.From)
		if err != nil {
			fmt.Println("Could not remove", seed.From, err)
		}
	}
}
