package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/tsaron/anansi"
	"tsaron.com/traccar-simulator/pkg/config"
	"tsaron.com/traccar-simulator/pkg/traccar"
)

func main() {
	var err error

	var env config.Env
	if err = anansi.LoadEnv(&env); err != nil {
		panic(err)
	}

	log := anansi.NewLogger(env.Name)

	// connect to db
	var db *pg.DB
	if db, err = config.SetupDB(env); err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Err(err).Msg("failed to disconnect from traccar DB cleanly")
		}
	}()
	log.Info().Msg("successfully connected to traccar DB")

	devices := traccar.ReadDeviceDump("samples/devices.csv")
	positions := traccar.ReadPositionDump("samples/positions.csv")

	if err := traccar.CreateDevices(db, devices); err != nil {
		panic(err)
	}

	for _, p := range positions {
		if err := traccar.CreatePosition(db, p); err != nil {
			panic(err)
		}
	}
}
