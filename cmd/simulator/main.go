package main

import (
	"os"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"
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

	var limit uint64
	limitStr := os.Args[1]
	if limitStr == "" {
		limit = 0
	} else {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			panic(errors.Wrap(err, "failed to parse limit"))
		}
	}

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

	if limit != 0 {
		positions = positions[:limit]
	}

	for _, p := range positions {
		if err := traccar.CreatePosition(db, p); err != nil {
			panic(err)
		}
	}
}
