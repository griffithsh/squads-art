package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"

	"github.com/griffithsh/squads-art/cmd/tiled-convert/tiled"
)

type flattenedTile struct {
	texture    string
	w, h, x, y int
	xOff, yOff int
}
type converter struct {
	maps     map[string]tiled.Map
	tilesets map[string]tiled.Tileset

	converted map[string]game.CombatMapRecipe
}

func newConverter() *converter {
	return &converter{
		maps:      map[string]tiled.Map{},
		tilesets:  map[string]tiled.Tileset{},
		converted: map[string]game.CombatMapRecipe{},
	}
}

func (c *converter) scan(dir string) error {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("abs dir %s: %v", dir, err)
	}
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %v", dir, err)
	}

	for _, info := range infos {
		if strings.HasSuffix(info.Name(), ".tmx.export") {
			f, err := os.Open(path.Join(absDir, info.Name()))
			if err != nil {
				return fmt.Errorf("open file: %v", err)
			}

			var v tiled.Map
			err = json.NewDecoder(f).Decode(&v)
			if err != nil {
				return fmt.Errorf("decode map: %v", err)
			}
			filename := strings.TrimSuffix(info.Name(), ".tmx.export")
			c.maps[filename] = v
		} else if strings.HasSuffix(info.Name(), ".tsx.export") {
			// handle
			f, err := os.Open(path.Join(absDir, info.Name()))
			if err != nil {
				return fmt.Errorf("open tileset: %v", err)
			}

			var v tiled.Tileset
			err = json.NewDecoder(f).Decode(&v)
			if err != nil {
				return fmt.Errorf("decode tileset: %v", err)
			}
			filename := strings.TrimSuffix(info.Name(), ".tsx.export")
			c.tilesets[filename] = v
		}
	}
	return nil
}

func getZIndex(props []tiled.Property) int {
	for _, prop := range props {
		if strings.ToUpper(prop.Name) == "Z" {
			switch v := prop.Value.(type) {
			case float64:
				return int(v)
			case string:
				i, err := strconv.Atoi(v)
				if err != nil {
					fmt.Printf("could not convert %s to an int for Z index\n", v)
				}
				return i
			case int:
				return v
			default:
				fmt.Printf("unhandled type %T for Z index", prop.Value)
				return 0
			}
		}
	}
	return 0
}
func (c *converter) convert() error {
	f := geom.NewField(34, 19, 40)
	for filename, m := range c.maps {
		// make a new entry in converted
		// TODO!

		if m.Orientation != "hexagonal" {
			return fmt.Errorf("non-hexagonal orientation %s: %s", m.Orientation, filename)
		}
		if m.TileWidth != 72 {
			return fmt.Errorf("incorrect tilewidth %d: %s", m.TileWidth, filename)
		}
		if m.TileHeight != 40 {
			return fmt.Errorf("incorrect tileheight %d: %s", m.TileHeight, filename)
		}
		if m.HexSideLength != 34 {
			return fmt.Errorf("incorrect hexsidelength %d: %s", m.HexSideLength, filename)
		}
		if m.StaggerAxis != "x" {
			return fmt.Errorf("incorrect staggeraxis %s: %s", m.StaggerAxis, filename)
		}
		if m.StaggerIndex != "odd" {
			return fmt.Errorf("incorrect staggerindex %s: %s", m.StaggerIndex, filename)
		}
		if m.RenderOrder != "right-up" {
			return fmt.Errorf("incorrect renderorder %s: %s", m.RenderOrder, filename)
		}
		// if tile render order is not Right Up...
		// return an err

		flatTiles := map[int]flattenedTile{}
		for _, mts := range m.Tilesets {
			tsx := c.tilesets[strings.TrimSuffix(mts.Source, ".tsx")]
			for _, t := range tsx.Tiles {
				flatTiles[t.ID+mts.FirstGID] = flattenedTile{
					w:    tsx.TileWidth,
					h:    tsx.TileHeight,
					x:    (t.ID % tsx.Columns) * tsx.TileWidth,
					y:    (t.ID / tsx.Columns) * tsx.TileHeight,
					xOff: tsx.TileOffset.X,
					yOff: tsx.TileOffset.Y,
					// FIXME: This is brittle!
					texture: path.Join("combat-terrain", tsx.Image),
				}
			}
		}
		hexes := []game.CombatMapRecipeHex{}
		starts := []geom.Key{}
		for _, layer := range m.Layers {
			for i, datum := range layer.Data {
				k := geom.Key{
					M: i % m.Width,
					N: i / m.Height,
				}
				flat := flatTiles[datum]
				hexes = append(hexes, game.CombatMapRecipeHex{
					Position: k,
					Obstacle: false, // FIXME
					Visuals: []game.CombatMapRecipeVisual{

						game.CombatMapRecipeVisual{
							Layer:   getZIndex(layer.Properties),
							XOffset: flat.xOff,
							YOffset: flat.yOff,
							Frames: []game.CombatMapRecipeHexFrame{
								game.CombatMapRecipeHexFrame{
									Texture: flat.texture,
									X:       flat.x,
									Y:       flat.y,
									W:       flat.w,
									H:       flat.h,
								},
							},
						},
					},
				})
			}

			for _, obj := range layer.Objects {
				if strings.ToLower(obj.Type) == "start" {
					k := f.Wtok(obj.X-float64(m.TileWidth/2), obj.Y-float64(m.TileHeight/2))
					starts = append(starts, k)
				}
			}
		}
		recipe := game.CombatMapRecipe{
			Hexes:  hexes,
			Starts: starts,
			TileW:  m.TileWidth,
			TileH:  m.TileHeight,
		}
		c.converted[filename+".terrain"] = recipe
	}
	return nil
}

func (c *converter) write(dir string) error {
	for k, r := range c.converted {
		b, err := json.Marshal(r)
		if err != nil {
			return fmt.Errorf("marshal: %v", err)
		}
		if err = ioutil.WriteFile(path.Join(dir, k), b, 0644); err != nil {
			return fmt.Errorf("write file: %v", err)
		}
	}
	return nil
}
