package mid

import (
	"context"

	"github.com/zaouldyeck/webservice/app/api/errs"
	"github.com/zaouldyeck/webservice/foundation/logger"
)

// Errors handles errors coming out of the call chain.
func Errors(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)
	if err == nil {
		return nil
	}

	log.Error(ctx, "message", "ERROR", err.Error())

	if errs.IsError(err) {
		return errs.GetError(err)
	}

	return errs.Newf(errs.Unknown, errs.Unknown.String())
}
