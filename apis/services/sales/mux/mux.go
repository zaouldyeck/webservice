// Package mid provides support to bind domain level routes
// to the application mid.
package mux

import (
	"os"

	"github.com/zaouldyeck/webservice/apis/services/api/mid"
	"github.com/zaouldyeck/webservice/apis/services/sales/route/sys/checkapi"
	"github.com/zaouldyeck/webservice/foundation/logger"
	"github.com/zaouldyeck/webservice/foundation/web"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log))

	checkapi.Routes(mux)

	return mux
}
