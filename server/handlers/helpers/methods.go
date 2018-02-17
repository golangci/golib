package helpers

import (
	"github.com/golangci/golib/server/context"
	"github.com/golangci/golib/server/handlers/herrors"
	"github.com/golangci/golib/server/handlers/manager"
)

func AllowMethods(h manager.Handler, methods []string) manager.Handler {
	return func(ctx context.C) error {
		methodOK := false
		for _, m := range methods {
			if ctx.R.Method == m {
				methodOK = true
				break
			}
		}

		if !methodOK {
			return herrors.New404Errorf("Unallowed method was used")
		}

		return h(ctx)
	}
}

func OnlyPOST(h manager.Handler) manager.Handler {
	return AllowMethods(h, []string{"POST"})
}

func OnlyPUT(h manager.Handler) manager.Handler {
	return AllowMethods(h, []string{"PUT"})
}

func OnlyDELETE(h manager.Handler) manager.Handler {
	return AllowMethods(h, []string{"DELETE"})
}

func OnlyGET(h manager.Handler) manager.Handler {
	return AllowMethods(h, []string{"GET"})
}
