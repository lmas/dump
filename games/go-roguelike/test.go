package main

import (
	"fmt"

	rl "github.com/lmas/go-rl/engine"
)

func main() {
	em := rl.NewEntityManager()
	id := em.NewEntity()
	fmt.Println("ent id:", id)
	em.All()

	fmt.Println("New compo")
	com := 2
	em.NewComponent(com, map[string]interface{}{
		"x": 1,
		"y": 1,
	})
	em.All()

	fmt.Println("Add com")
	em.AddComponent(id, 1)
	em.AddComponent(id, 1)
	em.AddComponent(id, 2)
	em.AddComponent(id, 3)
	em.All()

	fmt.Println("Del com")
	fmt.Println("comps for ent:", em.ComponentsForEntity(id))
	fmt.Println("ents with comp:", em.EntitiesWithComponent(com))
	em.DelComponent(id, 2)
	fmt.Println("comps for ent:", em.ComponentsForEntity(id))
	fmt.Println("ents with comp:", em.EntitiesWithComponent(com))
	em.All()

	fmt.Println("System runs")
	em.NewSystem(func(entities []int) {
		fmt.Println("Running system against:", entities)
	}, []int{1})
	em.RunSystems()
	em.All()

	fmt.Println("Removing ents")
	em.DelEntity(id)
	em.All()

}
