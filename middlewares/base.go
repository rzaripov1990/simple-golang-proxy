package middlewares

import (
	"context"
	"net/http"

	"github.com/rzaripov1990/simple-golang-proxy/utils"
	"github.com/rzaripov1990/simple-golang-proxy/utils/logger"
)

type key int

const ctxKey key = iota + 16

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			logRequestId   string = utils.RequestID()
			logRequestBody interface{}
			err            error
		)

		err = utils.JsonDecode(r, &logRequestBody, true)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.L(logRequestId, err)
			return
		}
		logger.L(logRequestId, logRequestBody)

		ctx := context.WithValue(r.Context(), ctxKey, logRequestId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

func GetRequestId(ctx context.Context) (value string) {
	val := ctx.Value(ctxKey)
	if val != nil {
		s, ok := val.(string)
		if ok {
			return s
		}
	}
	return
}
