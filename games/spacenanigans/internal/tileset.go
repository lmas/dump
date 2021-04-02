package internal

type Tile struct {
	Name     string `json:"name"`
	Blocking bool   `json:"blocking"`
	SheetX   int    `json:"sheetx"`
	SheetY   int    `json:"sheety"`
}

var defaultTileset = map[int]Tile{
	0: Tile{"", true, 0, 2},  // empty
	1: Tile{"", false, 1, 2}, // catwalk
	2: Tile{"", false, 2, 2}, // floor
	3: Tile{"", false, 3, 2}, // white floor
	4: Tile{"", false, 4, 2}, // red floor
	5: Tile{"", false, 5, 2}, // green floor
	6: Tile{"", false, 6, 2}, // blue floor
	7: Tile{"", false, 7, 2}, // yellow floor

	8: Tile{"", true, 0, 0}, // wall
	9: Tile{"", true, 0, 1}, // glass wall

	10: Tile{"", true, 0, 3},   // closed door
	11: Tile{"", false, 1, 3},  // open door
	12: Tile{"", true, 2, 3},   // closed secure door
	13: Tile{"", false, 3, 3},  // open secure door
	14: Tile{"", true, 4, 3},   // table
	15: Tile{"", true, 5, 3},   // closed closet
	16: Tile{"", true, 6, 3},   // open closet
	17: Tile{"", true, 7, 3},   // offline computer
	18: Tile{"", true, 8, 3},   // online computer
	19: Tile{"", true, 9, 3},   // closed crate
	20: Tile{"", true, 10, 3},  // open crate
	21: Tile{"", false, 11, 3}, // chair
	22: Tile{"", true, 12, 3},  // vending machine
	23: Tile{"", true, 13, 3},  // bed

	// these are bitmasked walls
	100: Tile{"", true, 0, 0},
	101: Tile{"", true, 1, 0},
	102: Tile{"", true, 2, 0},
	103: Tile{"", true, 3, 0},
	104: Tile{"", true, 4, 0},
	105: Tile{"", true, 5, 0},
	106: Tile{"", true, 6, 0},
	107: Tile{"", true, 7, 0},
	108: Tile{"", true, 8, 0},
	109: Tile{"", true, 9, 0},
	110: Tile{"", true, 10, 0},
	111: Tile{"", true, 11, 0},
	112: Tile{"", true, 12, 0},
	113: Tile{"", true, 13, 0},
	114: Tile{"", true, 14, 0},
	115: Tile{"", true, 15, 0},

	// bitmasked glass walls
	120: Tile{"", true, 0, 1},
	121: Tile{"", true, 1, 1},
	122: Tile{"", true, 2, 1},
	123: Tile{"", true, 3, 1},
	124: Tile{"", true, 4, 1},
	125: Tile{"", true, 5, 1},
	126: Tile{"", true, 6, 1},
	127: Tile{"", true, 7, 1},
	128: Tile{"", true, 8, 1},
	129: Tile{"", true, 9, 1},
	130: Tile{"", true, 10, 1},
	131: Tile{"", true, 11, 1},
	132: Tile{"", true, 12, 1},
	133: Tile{"", true, 13, 1},
	134: Tile{"", true, 14, 1},
	135: Tile{"", true, 15, 1},
}
