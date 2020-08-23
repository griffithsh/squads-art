package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	obstacle   game.ObstacleType
	obscures   tiled.ObscuresProperty
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

// getTileObstacle returns true when the tile is an obstacle.
func getTileObstacle(tileID int, tiles []tiled.TilesetTile) (game.ObstacleType, error) {
	for _, t := range tiles {
		if t.ID == tileID {
			for _, prop := range t.Properties {
				if strings.ToLower(prop.Name) == "obstacle" {
					// What enum value can we get from this string?
					str, ok := prop.Value.(string)
					if !ok {
						return game.NonObstacle, fmt.Errorf("non-string value for obstacle property of tile %d", t.ID)
					}
					parsed, err := game.ParseObstacleType(str)
					if err != nil {
						return game.NonObstacle, fmt.Errorf("parse %s: %v", str, err)
					}

					return parsed, nil
				}
			}
			break
		}
	}
	return game.NonObstacle, nil
}

func getObscuredOtherTiles(tileID int, tiles []tiled.TilesetTile) (tiled.ObscuresProperty, error) {
	for _, t := range tiles {
		if t.ID == tileID {
			for _, prop := range t.Properties {
				if strings.ToLower(prop.Name) == "obscures" {
					str, ok := prop.Value.(string)
					if !ok {
						fmt.Fprintf(os.Stderr, "tiledID %d: non string obscures property", tileID)
						return tiled.ObscuresProperty{}, nil
					}

					var result tiled.ObscuresProperty
					r := bufio.NewReader(bytes.NewReader([]byte(str)))
					for {
						line, err := r.ReadString('\n')
						if err == io.EOF {
							break
						} else if err != nil {
							return nil, fmt.Errorf("ReadString: %v", err)
						}
						vals := strings.Split(string(line), ",")
						if len(vals) < 2 {
							return nil, fmt.Errorf("split string: %v", err)
						}

						m, err := strconv.Atoi(strings.TrimSpace(vals[0]))
						if err != nil {
							return nil, fmt.Errorf("Atoi: %v", err)
						}
						n, err := strconv.Atoi(strings.TrimSpace(vals[1]))
						if err != nil {
							return nil, fmt.Errorf("Atoi: %v", err)
						}
						result = append(result, tiled.ObscuresPropertyCoordinate{M: m, N: n})
					}
					return result, nil
				}
			}
			break
		}
	}
	return tiled.ObscuresProperty{}, nil
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

		flatTiles := map[int]flattenedTile{}
		for _, mts := range m.Tilesets {
			tsx := c.tilesets[strings.TrimSuffix(mts.Source, ".tsx")]
			for i := 0; i < tsx.TileCount; i++ {
				obstacle, err := getTileObstacle(i, tsx.Tiles)
				if err != nil {
					panic(tsx.Name + ": obstacle: " + err.Error())
				}
				obscures, err := getObscuredOtherTiles(i, tsx.Tiles)
				if err != nil {
					panic(tsx.Name + ": obscured property: " + err.Error())
				}
				flatTiles[i+mts.FirstGID] = flattenedTile{
					w:    tsx.TileWidth,
					h:    tsx.TileHeight,
					x:    (i % tsx.Columns) * tsx.TileWidth,
					y:    (i / tsx.Columns) * tsx.TileHeight,
					xOff: tsx.TileOffset.X,
					yOff: tsx.TileOffset.Y,
					// FIXME: This is brittle!
					texture: path.Join("combat-terrain", tsx.Image),

					obstacle: obstacle,
					obscures: obscures,
				}
			}
		}
		hexes := []game.CombatMapRecipeHex{}
		starts := []geom.Key{}
		for _, layer := range m.Layers {
			for i, datum := range layer.Data {
				if datum == 0 {
					continue
				}
				k := geom.Key{
					M: i % m.Width,
					N: i / m.Height,
				}
				flat, ok := flatTiles[datum]
				if !ok {
					panic(fmt.Sprintf("no flattened tile for %d: %v", datum, flatTiles))
				}
				hexes = append(hexes, game.CombatMapRecipeHex{
					Position: k,
					Obstacle: flat.obstacle,
					Visuals: []game.CombatMapRecipeVisual{
						game.CombatMapRecipeVisual{
							Layer:    getZIndex(layer.Properties),
							XOffset:  flat.xOff,
							YOffset:  flat.yOff,
							Obscures: realise(k, flat.obscures),
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

func realise(p geom.Key, offsets tiled.ObscuresProperty) []geom.Key {
	result := make([]geom.Key, len(offsets))
	for i, off := range offsets {
		result[i] = geom.Key{M: p.M + off.M, N: p.N + off.N}
	}
	return result
}
