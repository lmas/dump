package engine

import "log"

// Item IDs
const (
	ITEM_WRENCH = iota
	ITEM_ORE
)

// Item types
const (
	ITEM_BLOCK = iota
	ITEM_TOOL
	ITEM_BUILDING
)

type Item struct {
	Name    string
	Display string
	ID      int
	Type    int
	Usef    func(*Item, *Mob) error
	Placef  func(*Item, *Mob) error
}

var ITEMS map[int]Item

func NewItem(id int) *Item {
	i := ITEMS[id]
	return &i
}

func LoadItems() {
	ITEMS = map[int]Item{
		ITEM_WRENCH: Item{"Wrench", "w", ITEM_WRENCH, ITEM_TOOL, wrench_use, nil},
		ITEM_ORE:    Item{"Ore", "o", ITEM_ORE, ITEM_BLOCK, nil, nil},
	}
}

func wrench_use(i *Item, m *Mob) error {
	t := MapGet(DirToMap(m))
	if t.Type >= TILE_HILL {
		if m.Block == nil && Prob(25) {
			log.Println("Gave a new ore block.")
			m.Block = NewItem(ITEM_ORE)
		}
	}
	return nil
}
