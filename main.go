package main

import (
	"fmt"
)

var (
  windowWidth float32 = 800
  windowHeight float32 = 360
  rows = 8
  cols = 14
)

type VelocityData struct {
  X, Y int
}

type SpriteData struct {}
type PosData struct {}

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

  Velocity.SetData(world, 1, VelocityData{5, 10})

  p1 := Velocity.Get(world, 1)

  p1.X /= 2
  p1.Y /= 2

  Velocity.Each(world, func(entity int, data *VelocityData) {
    data.X -= 1
    fmt.Println("entity, data", entity, data)
  })

  fmt.Println("data", world.stores[0].Data)
  fmt.Println("first", Velocity.First(world))
}

