package drink

import (
	"github.com/nicjohnson145/mixer-service/pkg/common"
	"database/sql"
	"net/http"
)

func createDrink(db *sql.DB) common.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
