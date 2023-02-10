package main

import (
	"flag"

	"github.com/nicjohnson145/mixer-service/pkg/auth"
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"github.com/nicjohnson145/mixer-service/pkg/db"
	"github.com/nicjohnson145/mixer-service/pkg/drink"
	"github.com/nicjohnson145/mixer-service/pkg/health"
	"github.com/nicjohnson145/mixer-service/pkg/settings"
	"github.com/nicjohnson145/mixer-service/pkg/slow"
	"github.com/nicjohnson145/mixer-service/pkg/static"
	"github.com/nicjohnson145/mixer-service/pkg/user"
	log "github.com/sirupsen/logrus"
)

var (
	slowDown bool
)

func init() {
	flag.BoolVar(&slowDown, "slow", false, "Introduce a sleep in every request. Useful for UI testing")
}

func main() {
	flag.Parse()

	app := common.NewApp()

	if slowDown {
		slow.SlowDown(app)
	}

	db := db.NewDBOrDie(common.DefaultedEnvVar("DB_PATH", "mixer.db"))

	if err := auth.Init(app, db); err != nil {
		log.Fatal(err)
	}

	if err := drink.Init(app, db); err != nil {
		log.Fatal(err)
	}

	if err := health.Init(app, db); err != nil {
		log.Fatal(err)
	}

	if err := settings.Init(app, db); err != nil {
		log.Fatal(err)
	}

	if err := user.Init(app, db); err != nil {
		log.Fatal(err)
	}

	if err := static.Init(app); err != nil {
		log.Fatal(err)
	}

	port := common.DefaultedEnvVar("PORT", "30000")

	log.Info("Listening on port ", port)
	app.Listen(":" + port)
}
