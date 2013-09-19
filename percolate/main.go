package main

func main() {
	m, _ := NewRandomMaterial(200, 200, 0.594)
	connected := Search(m)
	Draw(m, connected, "material.png")
}
