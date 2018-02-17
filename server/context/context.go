package context

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type C struct {
	W         http.ResponseWriter
	R         *http.Request
	Ctx       context.Context
	L         *log.Entry
	RequestID string
}

func (ctx C) URLVar(name string) string {
	return mux.Vars(ctx.R)[name]
}

func (ctx C) ReturnJSON(r interface{}) {
	b, err := json.Marshal(r)
	if err != nil {
		ctx.L.Errorf("can't marshal json: %s", err)
		return
	}

	ctx.W.Header().Add("Content-Type", "application/json; charset=UTF-8")
	ctx.W.WriteHeader(http.StatusOK)
	ctx.W.Write(b)
}

func (ctx C) RedirectTemp(toURL string) {
	http.Redirect(ctx.W, ctx.R, toURL, http.StatusTemporaryRedirect)
}
