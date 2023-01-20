package tiled

type MapObject struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Class string  `json:"class"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	// "width": 0,
	// "height": 0,
	// "point": true,
	// "rotation": 0,
	// "visible": true,
}

type MapLayer struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Data    []int       `json:"data"`
	Objects []MapObject `json:"objects"`
	Width   int         `json:"width"`
	Height  int         `json:"height"`
	X       int         `json:"x"`
	Y       int         `json:"y"`
	//   "opacity":1,
	//   "visible":true,
	Properties []Property `json:"properties"`
}

type Property struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type Map struct {
	Infinite      bool         `json:"infinite"`
	Width         int          `json:"width"`
	Height        int          `json:"height"`
	HexSideLength int          `json:"hexsidelength"`
	Layers        []MapLayer   `json:"layers"`
	Orientation   string       `json:"orientation"`
	RenderOrder   string       `json:"renderorder"`
	StaggerAxis   string       `json:"staggeraxis"`
	StaggerIndex  string       `json:"staggerindex"`
	TileWidth     int          `json:"tilewidth"`
	TileHeight    int          `json:"tileheight"`
	Tilesets      []MapTileset `json:"tilesets"`
	Properties    []Property   `json:"properties"`
	//  "nextlayerid":2,
	//  "nextobjectid":1,
	// "compressionlevel":-1,
	//  "type":"map",
	//  "version":1.2,
	//  "tiledversion":"1.3.3",
}

// ObscuresProperty is a type that the "Obscures" custom property can be
// unmarshaled into.
type ObscuresProperty []ObscuresPropertyCoordinate

// ObscuresPropertyCoordinate is a hexagonal offset from the origin of the
// placed tile that defines an other hexagon that is obscured yb this tile.
type ObscuresPropertyCoordinate struct {
	M, N int
}

type MapTileset struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}

type TilesetTile struct {
	ID         int        `json:"id"`
	Properties []Property `json:"properties"`
}

type TilesetOffset struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Tileset struct {
	Columns int    `json:"columns"`
	Image   string `json:"image"`
	// "imageheight": 128,
	// "imagewidth": 128,
	// "margin": 0,
	Name string `json:"name"`
	// "spacing": 0,
	TileCount int `json:"tilecount"`
	// "tiledversion": "1.3.3",
	TileHeight int `json:"tileheight"`
	TileWidth  int `json:"tilewidth"`
	// "type": "tileset",
	// "version": 1.2
	Tiles      []TilesetTile `json:"tiles"`
	TileOffset TilesetOffset `json:"tileoffset"`
}
