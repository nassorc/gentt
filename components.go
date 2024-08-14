package main

func Vec2Add(a, b Vec2) Vec2 {
  return Vec2{a.X + b.X, a.Y + b.Y}
}

func Vec2Subtract(a, b Vec2) Vec2 {
  return Vec2{a.X - b.X, a.Y - b.Y}
}

func Vec2Scale(a Vec2, s float32) Vec2 {
  return Vec2{a.X * s, a.Y * s}
}

type Vec2 struct { X, Y float32 }

type EntityId = int

type Pos Vec2
