package static

import (
	"embed"
	_ "embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed public
var publicWebUiContent embed.FS

func Init(app *fiber.App) error {
	addStaticRoutes(app)
	return nil
}

func addStaticRoutes(app *fiber.App) {
	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(publicWebUiContent),
		PathPrefix: "public/webui",
		Browse:     true,
	}))
}
