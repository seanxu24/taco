package handler

import (
	"net/http"

	"book/service/search/cmd/api/internal/logic"
	"book/service/search/cmd/api/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func pingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPingLogic(r.Context(), svcCtx)
		err := l.Ping()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
