package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/griffithsh/squads/game"

	"github.com/griffithsh/squads-art/cmd/tiled-convert/tiled"
)

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

func (c *converter) convert() error {
	for filename, m := range c.maps {
		// make a new entry in converted
		// TODO!

		if m.Orientation != "hexagonal" {
			return fmt.Errorf("non-hexagonal orientation %s: %s", m.Orientation, filename)
		}
		// if tile width is not 72 ...
		// if tile height is not 40 ...
		// if tile side length is not 34...
		// if stagger axis is not X...
		// if stagger index is not Odd...
		// if tile render order is not Right Up...
		// return an err

		recipe := game.CombatMapRecipe{
			// Hexes: ...,
			// Starts: ...,
			TileW: m.TileWidth,
			TileH: m.TileHeight,
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
