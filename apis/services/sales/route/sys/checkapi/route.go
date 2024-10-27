package checkapi

import (
	"github.com/zaouldyeck/webservice/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
	app.HandleFunc("GET /testerror", testError)
}
