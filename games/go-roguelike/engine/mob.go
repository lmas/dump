package engine

type Mob struct {
	Name    string
	Display string
	Pos     *Vector
}

var MobList []*Mob

func NewMob(name string, display string, x, y int) Mob {
	mob := Mob{name, display, NewVector(float64(x), float64(y))}
	MobList = append(MobList, &mob)
	mob.Move(0, 0) // Updates the map cell so it contains the mob
	return mob
}

func (m *Mob) Move(x, y int) {
	delta := NewVector(float64(x), float64(y))
	newpos := m.Pos.Add(delta)
	if MapCollisionCheck(newpos) == false {
		return
	}

	// Set news pos
	oldpos := m.Pos
	m.Pos = newpos

	// Update map cells
	oldcell := MapGet(int(oldpos.X), int(oldpos.Y))
	newcell := MapGet(int(m.Pos.X), int(m.Pos.Y))
	oldcell.Mob = nil
	newcell.Mob = m
	MapSet(int(oldpos.X), int(oldpos.Y), oldcell)
	MapSet(int(m.Pos.X), int(m.Pos.Y), newcell)
}

func (m *Mob) Life() {
	m.RandomWalk()
}

func (m *Mob) RandomWalk() {
	//Randomly walk around
	if Prob(50) {
		m.Move(RandRange(-1, 1), 0)
	} else {
		m.Move(0, RandRange(-1, 1))
	}
}

func (m *Mob) Track(t *Mob) {
	dist := m.Pos.Dist(t.Pos)

	if dist < 2 {
		// Don't get too close
		m.StepTo(t.Pos)
	} else if dist > 5 {
		// Don\t get too far away
		m.StepFrom(t.Pos)
	} else {
		m.RandomWalk()
	}
}

func (m *Mob) StepTo(v *Vector) {
	delta := m.Pos.Sub(v)
	if Prob(50) {
		if delta.X > 0 {
			m.Move(1, 0)
		} else if delta.X < 0 {
			m.Move(-1, 0)
		} else if delta.Y > 0 {
			m.Move(0, 1)
		} else if delta.Y < 0 {
			m.Move(0, -1)
		}
	} else {
		if delta.Y > 0 {
			m.Move(0, 1)
		} else if delta.Y < 0 {
			m.Move(0, -1)
		} else if delta.X > 0 {
			m.Move(1, 0)
		} else if delta.X < 0 {
			m.Move(-1, 0)
		}
	}
}

func (m *Mob) StepFrom(v *Vector) {
	delta := m.Pos.Sub(v)
	if Prob(50) {
		if delta.X > 0 {
			m.Move(-1, 0)
		} else if delta.X < 0 {
			m.Move(1, 0)
		} else if delta.Y > 0 {
			m.Move(0, -1)
		} else if delta.Y < 0 {
			m.Move(0, 1)
		}
	} else {
		if delta.Y > 0 {
			m.Move(0, -1)
		} else if delta.Y < 0 {
			m.Move(0, 1)
		} else if delta.X > 0 {
			m.Move(-1, 0)
		} else if delta.X < 0 {
			m.Move(1, 0)
		}
	}
}
