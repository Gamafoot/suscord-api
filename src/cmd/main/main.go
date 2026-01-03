package main

import (
	"fmt"
	"suscord/internal/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	if err = app.Run(); err != nil {
		panic(err)
	}
}
