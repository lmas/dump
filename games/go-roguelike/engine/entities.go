package engine

import "fmt"

type EntitySystem struct {
	UpdateFunc         func([]int)
	RequiredComponents []int
}

type EntityManager struct {
	latestid int
	// list of all available systems that should be run
	systems []*EntitySystem
	// key = component id, value = map of new data values to add to entity
	components map[int]map[string]interface{}
	// key = entity id, value = list of component ids
	entities map[int][]int
	// key = entity id, value = map of data for entity
	entitydata map[int]map[string]interface{}
}

func NewEntityManager() *EntityManager {
	em := EntityManager{
		latestid:   -1,
		systems:    *new([]*EntitySystem),
		components: make(map[int]map[string]interface{}),
		entities:   make(map[int][]int),
		entitydata: make(map[int]map[string]interface{}),
	}
	return &em
}

// TODO: TEST FUNC TO BE REMOVED
func (em *EntityManager) All() {
	fmt.Println("entitites:\t", em.entities)
	fmt.Println("entdata:\t", em.entitydata)
	fmt.Println("components:\t", em.components)
	fmt.Println()
}

// Create a new entity and return it's ID
func (em *EntityManager) NewEntity() int {
	em.latestid++
	id := em.latestid
	em.entities[id] = *new([]int)
	em.entitydata[id] = make(map[string]interface{})
	return id
}

// Remove existing entity, using it's ID
func (em *EntityManager) DelEntity(id int) {
	delete(em.entities, id)
	delete(em.entitydata, id)
}

// Add component to an entity
func (em *EntityManager) AddComponent(id, component int) {
	// Check if the component is already in the list
	for _, val := range em.entities[id] {
		if val == component {
			// It's already here, let's bail out
			return
		}
	}
	// If not, we're good to add it now
	em.entities[id] = append(em.entities[id], component)

	// Add any new data from the component
	for key, val := range em.components[component] {
		if _, exists := em.entitydata[id][key]; exists == true {
			continue
		}
		em.entitydata[id][key] = val
	}
}

// Remove a component from an entity
func (em *EntityManager) DelComponent(id, component int) {
	for i, val := range em.entities[id] {
		if val == component {
			// It can get kinda ugly to remove values from slices in go...
			// Once we've found the component in the list we have to
			// split up the slice in two and then smash them together
			// again, but without the component
			em.entities[id] = append(
				em.entities[id][:i],
				em.entities[id][i+1:]...,
			)

			// Remove data values related to the component
			for key := range em.components[component] {
				delete(em.entitydata[id], key)
			}
			return
		}
	}
}

// Return a list of components for an entity
func (em *EntityManager) ComponentsForEntity(id int) []int {
	return em.entities[id]
}

// Return a list of entities with a component
func (em *EntityManager) EntitiesWithComponent(cid int) []int {
	ents := *new([]int)
	for id, coms := range em.entities {
		for _, c := range coms {
			if c != cid {
				continue
			}
			ents = append(ents, id)
		}
	}
	return ents
}

// Add a new component to the system
func (em *EntityManager) NewComponent(cid int, data map[string]interface{}) {
	em.components[cid] = data
}

// Add a new system to be run against entities with matching components
func (em *EntityManager) NewSystem(f func([]int), coms []int) {
	s := EntitySystem{
		UpdateFunc:         f,
		RequiredComponents: coms,
	}
	em.systems = append(em.systems, &s)
}

// Run all available systems once
func (em *EntityManager) RunSystems() {
	// TODO: Store a prepared list of entities instead of doing a loop each time
	// Should also enforce the entities to have all the required components too?
	for _, s := range em.systems {
		// First get a list of entities to run the system against
		ents := *new([]int)
		for _, c := range s.RequiredComponents {
			ents = append(ents, em.EntitiesWithComponent(c)...)
		}
		// Then run the system on this new list
		s.UpdateFunc(ents)
	}
}
