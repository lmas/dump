package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var inFile string
	var outFile string
	if len(os.Args) > 1 {
		inFile = os.Args[1]
	}
	if inFile == "" {
		panic(fmt.Errorf("missing map file path parameter"))
	}

	if len(os.Args) > 2 {
		outFile = os.Args[2]
	}
	if outFile == "" {
		ext := filepath.Ext(inFile)
		base := filepath.Base(inFile)
		outFile = strings.TrimSuffix(base, ext) + ".txt"
	}

	b, err := loadByondMap(inFile)
	if err != nil {
		panic(err)
	}

	m, err := convertMap(b)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(m)
	if err != nil {
		panic(err)
	}
}

////////////////////////////////////////////////////////////////////////////////

var (
	rSpawn = regexp.MustCompile(`/landmark/start{name = "([^"]*)"}`)
)

type ParsedMap struct {
	Width  int   `json:"width"`
	Height int   `json:"height"`
	SpawnX int   `json:"spawnx"`
	SpawnY int   `json:"spawny"`
	Raw    []int `json:"raw"`
}

var ignore = []string{
	"mineral",
	"beach",
	"shuttle",
	"space",
	"snow",
	"roads",
}

func convertMap(b *ByondMap) (*ParsedMap, error) {
	m := &ParsedMap{
		Width:  b.Width,
		Height: b.Height,
	}

	i := 0
	for _, k := range b.Raw {
		tile, found := b.Tiles[k]
		if !found {
			log.Printf("Unknown tile: %s\n", k)
			continue
		}

		switch {
		//case tile.Turf == "space":
		//m.Raw = append(m.Raw, 0)

		case strings.Contains(tile.Def, "computer"):
			m.Raw = append(m.Raw, 18)

		case strings.Contains(tile.Def, "table"):
			m.Raw = append(m.Raw, 14)

		case strings.Contains(tile.Def, "closet/crate"):
			m.Raw = append(m.Raw, 19)
		case strings.Contains(tile.Def, "closet"):
			m.Raw = append(m.Raw, 15)

		case strings.Contains(tile.Def, "chair"):
			m.Raw = append(m.Raw, 21)

		case strings.Contains(tile.Def, "structure/bed"):
			m.Raw = append(m.Raw, 23)

		case strings.Contains(tile.Def, "vending"):
			m.Raw = append(m.Raw, 22)

		case strings.Contains(tile.Def, "door/airlock/external"):
			m.Raw = append(m.Raw, 12)
		case strings.Contains(tile.Def, "door/airlock"):
			m.Raw = append(m.Raw, 11)
		case strings.Contains(tile.Def, "door/window"):
			m.Raw = append(m.Raw, 11)
		case strings.Contains(tile.Def, "door/unpowered"):
			m.Raw = append(m.Raw, 10)

		case strings.Contains(tile.Turf, "wall"):
			m.Raw = append(m.Raw, 8)
		case strings.Contains(tile.Def, "/obj/structure/window"):
			m.Raw = append(m.Raw, 9)

		case strings.Contains(tile.Def, "catwalk"):
			m.Raw = append(m.Raw, 1)
		case strings.Contains(tile.Turf, "floor") || strings.Contains(tile.Turf, "pool"):
			m.Raw = append(m.Raw, 2)

		default:
			found = false
			for _, s := range ignore {
				if strings.Contains(tile.Turf, s) {
					found = true
					break
				}
			}
			m.Raw = append(m.Raw, 0)
			if !found {
				log.Printf("Unhandled turf: %s\n", tile.Turf)
			}
		}

		i += 1
		spawn := rSpawn.FindStringSubmatch(tile.Def)
		if len(spawn) == 2 {
			//fmt.Println(spawn[1])
			if spawn[1] == "Captain" {
				m.SpawnX, m.SpawnY = i%b.Width, i/b.Width
			}
		}
	}
	// Ensure the spawn point is not blocked
	m.Raw[m.SpawnY*b.Width+m.SpawnX] = 2

	return m, nil
}

////////////////////////////////////////////////////////////////////////////////

var (
	rDef  = regexp.MustCompile(`^\"([a-zA-Z]+)\" = \((.+)\)$`)
	rTurf = regexp.MustCompile(`/turf/([a-z0-9_\-\/]*)`)
	rRow  = regexp.MustCompile(`^([a-zA-Z]+)$`)
)

type ByondTile struct {
	Key  string
	Def  string
	Turf string
}

type ByondMap struct {
	Width  int
	Height int
	Tiles  map[string]*ByondTile
	Raw    []string
}

func loadByondMap(path string) (*ByondMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := &ByondMap{
		Tiles: make(map[string]*ByondTile),
	}
	keyLen := 0
	scanner := bufio.NewScanner(f) // defaults to line scanning
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		tile := parseTile(line)
		if tile != nil {
			m.Tiles[tile.Key] = tile
			if keyLen == 0 {
				keyLen = len(tile.Key)
			}
			continue
		}

		row := splitRow(keyLen, line)
		if row != nil {
			if m.Width == 0 {
				m.Width = len(row)
			}
			m.Height += 1
			for _, r := range row {
				t, found := m.Tiles[r]
				if found {
					m.Raw = append(m.Raw, t.Key)
				}
			}
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func parseTile(line string) *ByondTile {
	d := rDef.FindStringSubmatch(line)
	if len(d) != 3 {
		return nil
	}
	key, def := d[1], d[2]
	t := rTurf.FindAllStringSubmatch(def, -1)
	if len(t) < 1 {
		return nil
	}
	turf := t[len(t)-1][1] // assume the last turf in the list is the one we want
	b := &ByondTile{
		Key:  key,
		Def:  def,
		Turf: turf,
	}
	return b
}

func splitRow(keyLen int, line string) []string {
	r := rRow.FindString(line)
	if r == "" {
		return nil
	}
	var row []string
	for i := 0; i < len(r)/keyLen; i++ {
		k := r[i*keyLen : (i*keyLen)+keyLen]
		row = append(row, k)
	}
	return row
}
