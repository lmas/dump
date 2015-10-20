package engine

import (
	"fmt"
	"math/rand"
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
}

type GameMap struct {
	Tiles *map[int]MapTile
	Cells *map[string]*MapCell
}

var CurrentMap GameMap

func LoadTileSet() *map[int]MapTile {
	tiles := make(map[int]MapTile)
	tiles[TILE_PLAIN] = MapTile{"Plain", ".", B_GREEN, false}
	tiles[TILE_HILL] = MapTile{"Hill", "^", B_GREEN, false}
	tiles[TILE_FOREST] = MapTile{"Forest", "#", B_GREEN, false}
	tiles[TILE_SAND] = MapTile{"Sand", ",", B_YELLOW, false}
	tiles[TILE_WATER] = MapTile{"Water", "~", B_BLUE, true}
	tiles[TILE_MOUNTAIN] = MapTile{"Mountain", "^", B_WHITE, true}
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

func NewCell(x, y int, tile int) *MapCell {
	return &MapCell{tile, NewVector(float64(x), float64(y)), nil}
}

func GenerateCell(x, y int) *MapCell {
	noise := Noise(x, y)

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
	cell := NewCell(x, y, tile)
	MapSet(x, y, cell)
	return cell
}

func PosToMap(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func MapSet(x, y int, cell *MapCell) {
	(*CurrentMap.Cells)[PosToMap(x, y)] = cell
}

func MapGet(x, y int) *MapCell {
	// TODO: Check for missing cells
	cell := (*CurrentMap.Cells)[PosToMap(x, y)]
	if cell == nil {
		cell = GenerateCell(x, y)
	}
	return cell
}

func GetMapScreenSize() (x, y int) {
	x, y = GetScreenSize()
	return x - 20, y - 2
}

func MapCollisionCheck(v *Vector) bool {
	cell := MapGet(int(v.X), int(v.Y))
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
			cell := MapGet(tmpx, tmpy)
			if cell == nil {
				ScreenPrint(x, y, RED, "!")
				continue
			}

			if cell.Mob != nil {
				ScreenPrint(x, y, B_WHITE, cell.Mob.Display)
			} else if cell.Type > 0 {
				tile := TileGet(cell.Type)
				ScreenPrint(x, y, tile.Color, tile.Display)
			} else {
				ScreenPrint(x, y, WHITE, " ")
			}
		}
	}
}
