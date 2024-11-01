package mid

import (
	"context"
	"net/http"

	"github.com/zaouldyeck/webservice/app/api/mid"
	"github.com/zaouldyeck/webservice/foundation/web"
)

// Metrics updates program counters using the middleware functionality.
func Metrics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Metrics(ctx, hdl)
		}

		return h
	}

	return m
}
