package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

const (
	cMaxSize int = 5000
)

type Map struct {
	Width   int          `json:"width"`
	Height  int          `json:"height"`
	SpawnX  int          `json:"spawnx"`
	SpawnY  int          `json:"spawny"`
	Tileset map[int]Tile `json:"tileset"`
	Raw     []int        `json:"raw"`
}

func LoadMap(path string) (*Map, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var m *Map
	d := json.NewDecoder(f)
	err = d.Decode(&m)
	if err != nil {
		return nil, err
	}
	m.Tileset = defaultTileset

	if m.Width < 1 || m.Width > cMaxSize {
		return nil, fmt.Errorf("invalid map width: %d", m.Width)
	}
	if m.Height < 1 || m.Height > cMaxSize {
		return nil, fmt.Errorf("invalid map height: %d", m.Height)
	}
	if m.SpawnX < 1 || m.SpawnX > m.Width {
		return nil, fmt.Errorf("invalid map spawn X: %v", m.SpawnX)
	}
	if m.SpawnY < 1 || m.SpawnY > m.Height {
		return nil, fmt.Errorf("invalid map spawn Y: %v", m.SpawnY)
	}
	if len(m.Raw) != m.Width*m.Height {
		return nil, fmt.Errorf("invalid map size: %d, expected %d", len(m.Raw), m.Width*m.Height)
	}
	return m, nil
}

func (m *Map) GetTile(pos Vector) int {
	x, y := int(math.Round(pos.X)), int(math.Round(pos.Y))
	// Any positions outside the map dimensions will be returned as tile=0 (space, blocking)
	if x <= 0 || x >= m.Width-1 || y <= 0 || y >= m.Height-1 {
		return 0
	}
	return m.Raw[y*m.Width+x]
}

func (m *Map) IsBlocking(pos Vector) bool {
	t := m.GetTile(pos)
	tile, found := m.Tileset[t]
	if !found || tile.Blocking {
		return true
	}
	return false
}
