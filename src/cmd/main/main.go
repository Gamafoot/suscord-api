package main

import (
	"log"
	"suscord/internal/app"
	"time"

	"github.com/avast/retry-go/v4"
)

func main() {
	var (
		a   *app.App
		err error
	)

	err = retry.Do(func() error {
		a, err = app.NewApp()
		if err != nil {
			log.Fatalf("%+v\n", err)
		}
		return err
	}, retry.Attempts(3), retry.Delay(time.Second))
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	if err = a.Run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
