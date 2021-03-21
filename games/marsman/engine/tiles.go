package engine

import "math/rand"

const (
	DIR_NORTH = iota
	DIR_EAST
	DIR_SOUTH
	DIR_WEST
)

const (
	TILE_ERROR = iota
	TILE_PLAIN
	TILE_HILL
	TILE_FOREST
	TILE_SAND
	TILE_WATER
	TILE_MOUNTAIN
	TILE_FLOOR
	TILE_WALL
)

type MapTile struct {
	Name          string
	Display       string
	Color         int16
	BlockMovement bool
}

type MapCell struct {
	Type int
	Pos  *Vector
	Mob  *Mob
	Item *Item
}

type GameMap struct {
	Tiles *map[int]MapTile
	Cells *map[string]*MapCell
}

var CurrentMap GameMap

func LoadTileSet() *map[int]MapTile {
	tiles := make(map[int]MapTile)
	tiles[TILE_PLAIN] = MapTile{"Plain", ".", B_YELLOW, false}
	tiles[TILE_HILL] = MapTile{"Hill", "^", B_GREEN, false}
	tiles[TILE_FOREST] = MapTile{"Forest", "#", B_YELLOW, false}
	tiles[TILE_SAND] = MapTile{"Sand", ",", YELLOW, false}
	tiles[TILE_WATER] = MapTile{"Water", "~", B_RED, true}
	tiles[TILE_MOUNTAIN] = MapTile{"Mountain", "^", B_BLACK, true}
	tiles[TILE_FLOOR] = MapTile{"Floor", "_", B_WHITE, false}
	tiles[TILE_WALL] = MapTile{"Wall", "#", B_WHITE, true}
	return &tiles
}

func TileGet(tile int) MapTile {
	return (*CurrentMap.Tiles)[tile]
}

func LoadMap() {
	tiles := LoadTileSet()
	cells := make(map[string]*MapCell)
	CurrentMap = GameMap{tiles, &cells}
	//for x := 0; x < 80; x++ {
	//for y := 0; y < 24; y++ {
	//MapSet(x, y, NewCell(x, y, TILE_PLAIN))
	//}
	//}
	//MapSet(0, 0, NewCell(0, 0, TILE_FOREST))
}

func NewCell(v *Vector, tile int) *MapCell {
	return &MapCell{tile, v, nil, nil}
}

func GenerateCell(v *Vector) *MapCell {
	noise := Noise(v.X, v.Y)

	tile := 0
	switch {
	case noise < -0.6:
		tile = TILE_WATER
	case noise > -0.6 && noise < -0.5:
		tile = TILE_SAND
	case noise > 0.3 && noise <= 0.7:
		tile = TILE_FOREST
	case noise > 0.7:
		tile = TILE_MOUNTAIN
	default:
		if rand.Float32() < 0.01 {
			tile = TILE_HILL
		} else {
			tile = TILE_PLAIN
		}
	}
	cell := NewCell(v, tile)
	MapSet(v, cell)
	return cell
}

func MapSet(v *Vector, cell *MapCell) {
	(*CurrentMap.Cells)[v.Str()] = cell
}

func MapGet(v *Vector) *MapCell {
	// TODO: Check for missing cells
	cell := (*CurrentMap.Cells)[v.Str()]
	if cell == nil {
		cell = GenerateCell(v)
	}
	return cell
}

func DirToMap(m *Mob) *Vector {
	dirs := map[int]*Vector{
		DIR_NORTH: NewVector(0, -1),
		DIR_EAST:  NewVector(1, 0),
		DIR_SOUTH: NewVector(0, 1),
		DIR_WEST:  NewVector(-1, 0),
	}
	dir := dirs[m.Dir]
	return m.Pos.Add(dir)
}

func GetMapScreenSize() (x, y int) {
	x, y = GetScreenSize()
	return x - 20, y - 2
}

func MapCollisionCheck(v *Vector) bool {
	cell := MapGet(v)
	tile := TileGet(cell.Type)
	return !tile.BlockMovement
}

func MapDraw(v *Vector, rad int) {
	maxx, maxy := GetMapScreenSize()
	halfx := maxx / 2
	halfy := maxy / 2

	for y := 0; y < maxy; y++ {
		for x := 0; x < maxx; x++ {
			tmpx := x + int(v.X) - halfx
			tmpy := y + int(v.Y) - halfy
			cell := MapGet(NewVector(float64(tmpx), float64(tmpy)))
			if cell == nil {
				ScreenPrint(x, y, RED, "!")
				continue
			}

			if cell.Mob != nil {
				ScreenPrint(x, y, B_WHITE, cell.Mob.Display)
			} else if cell.Item != nil {
				ScreenPrint(x, y, B_WHITE, cell.Item.Display)
			} else if cell.Type > 0 {
				tile := TileGet(cell.Type)
				ScreenPrint(x, y, tile.Color, tile.Display)
			} else {
				ScreenPrint(x, y, WHITE, " ")
			}
		}
	}
}
