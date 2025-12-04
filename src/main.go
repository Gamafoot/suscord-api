package main

import (
	"fmt"
	"mime"
)

func main() {
	a := mime.TypeByExtension(".png")
	fmt.Println(a)
}
