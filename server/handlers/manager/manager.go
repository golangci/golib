package manager

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	cntxt "context"

	"github.com/golangci/golib/server/context"
	"github.com/golangci/golib/server/handlers/herrors"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type Handler func(ctx context.C) error

type routerCB func(r *mux.Router)

var registry []routerCB

var logger = log.New()

func init() {
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.TextFormatter{ForceColors: true}
}

func Register(match string, handler Handler) {
	registry = append(registry, func(r *mux.Router) {
		r.HandleFunc(match, wrap(handler))
	})
}

func RegisterCallback(cb func(r *mux.Router)) {
	registry = append(registry, cb)
}

func wrap(f Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := fmt.Sprintf("0x%x", rand.Uint64())
		ctx := context.C{
			R:         r,
			W:         w,
			Ctx:       cntxt.Background(),
			RequestID: reqID,
			L:         logger.WithField("req_id", reqID),
		}
		startedAt := time.Now()
		err := f(ctx)
		processHandlerError(ctx, startedAt, err)
	}
}

func processHandlerError(ctx context.C, startedAt time.Time, err error) {
	if err == nil {
		ctx.L.Infof("%s[%s] successfully handled request for %s", ctx.R.Method, ctx.R.RequestURI, time.Since(startedAt))
		return
	}

	var code int
	if he, ok := err.(herrors.HTTPError); ok {
		code = he.Code()
	} else {
		code = http.StatusInternalServerError
	}
	ctx.L.Errorf("%s[%s] return %d: error during request processing: %s", ctx.R.Method, ctx.R.RequestURI, code, err)
	ctx.W.WriteHeader(code)
}

func MountHandlers(r *mux.Router) {
	for _, cb := range registry {
		cb(r)
	}
}

func GetHTTPHandler() http.Handler {
	r := mux.NewRouter()
	MountHandlers(r)
	return r
}
