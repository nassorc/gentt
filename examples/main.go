package main

import (
	"fmt"
  ecs "gent"
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

var Velocity = ecs.CreateComponent[VelocityData]()
var Sprite = ecs.CreateComponent[SpriteData]()
var Position = ecs.CreateComponent[PosData]()

type Val struct { 
  val int
}

func main() {
  world := ecs.NewWorld()
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

  fmt.Println("first", Velocity.First(world))

  for _, entity := range world.Query(Velocity, Sprite) {
    vel := Velocity.Get(world, entity)
    fmt.Println("query velocity and sprite", entity, vel)
  }
}

