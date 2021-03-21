package engine

import (
	"fmt"
	"log"
)

type Mob struct {
	Name    string
	Display string
	Pos     *Vector
	Dir     int
	Tool    *Item
	Block   *Item
}

func NewMob(name, display string, x, y int) *Mob {
	mob := &Mob{name, display, NewVector(float64(x), float64(y)), DIR_NORTH, nil, nil}
	mob.Move(x, y) // Updates the map cell so it contains the mob
	return mob
}

func (m *Mob) Move(x, y int) {
	delta := NewVector(float64(x), float64(y))
	newpos := m.Pos.Add(delta)
	if MapCollisionCheck(newpos) == false {
		return
	}

	// Set new pos
	oldpos := m.Pos
	m.Pos = newpos

	// Update map cells
	oldcell := MapGet(oldpos)
	oldcell.Mob = nil
	MapSet(oldpos, oldcell)

	newcell := MapGet(m.Pos)
	newcell.Mob = m
	MapSet(m.Pos, newcell)
}

func (m *Mob) Use() error {
	pos := DirToMap(m)
	c := MapGet(pos)

	// TODO
	if c.Item != nil {
		//switch c.Item.Type {
		//case ITEM_BLOCK:
		//if m.Block == nil {
		//m.Block = c.Item
		//c.Item = nil
		//MapSet(pos, c)
		//}
		//case ITEM_TOOL:
		//if m.Tool == nil {
		//m.Tool = c.Item
		//c.Item = nil
		//MapSet(pos, c)
		//}
		//}

		var item *Item
		if c.Item.Type == ITEM_BLOCK && m.Block == nil {
			item = c.Item
			m.Block = item
		} else if c.Item.Type == ITEM_TOOL && m.Tool == nil {
			item = c.Item
			m.Tool = item
		}
		if item != nil {
			c.Item = nil
			MapSet(pos, c)
			log.Printf("Grabbed %s at %s\n", item.Name, pos.Str())
			return nil
		}
	} else if m.Tool != nil && m.Tool.Usef != nil {
		return m.Tool.Usef(m.Tool, m)
	}
	return fmt.Errorf("Nothing to use")
}

func (m *Mob) Place() error {
	pos := DirToMap(m)
	c := MapGet(pos)
	if c.Item != nil {
		return fmt.Errorf("Already an item at %s", pos.Str())
	}

	var item *Item
	if m.Block != nil {
		item = m.Block
		m.Block = nil
	} else if m.Tool != nil {
		item = m.Tool
		m.Tool = nil
	} else {
		return fmt.Errorf("Nothing to place.")
	}

	c.Item = item
	MapSet(pos, c)
	if item.Placef != nil {
		return item.Placef(item, m)
	}
	log.Printf("Dropped %s at %s\n", item.Name, pos.Str())
	return nil
}
