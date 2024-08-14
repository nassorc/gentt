package main

import (
	"fmt"
)

var (
	windowWidth  float32 = 800
	windowHeight float32 = 360
	rows                 = 8
	cols                 = 14
)

type VelocityData struct {
	X, Y int
}

type SpriteData struct{}
type PosData struct{}

var Velocity = CreateComponent[VelocityData]()
var Sprite = CreateComponent[SpriteData]()
var Position = CreateComponent[PosData]()

type Val struct {
	val int
}

func main() {
	world := NewWorld()
	world.Create(Velocity)
	world.Create(Velocity)
	world.Create(Velocity, Sprite)

	Velocity.SetData(world, 1, VelocityData{5, 10})
	Velocity.SetData(world, 1, VelocityData{100, 100})

	// Velocity.SetData(world, 1, VelocityData{20, 21})
	// p1 := Velocity.Get(world, 1)
	//
	// p1.X /= 2
	// p1.Y /= 2

	Velocity.Each(world, func(entity int, data *VelocityData) {
		// data.X -= 1
		fmt.Println("entity, data", entity, data)
	})

	fmt.Println("data", world.stores[0].Dense)
	fmt.Println("first", Velocity.First(world))

	for _, entity := range world.Query(Velocity, Sprite) {
		vel := Velocity.Get(world, entity)
		fmt.Println("query velocity and sprite", entity, vel)
	}
}
