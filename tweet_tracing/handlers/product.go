package handlers

import (
	"../data"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
)

type TweetTracingHandler struct {
	Tracer trace.Tracer
	Repo   data.TweetTracingRepository
}

func (h TweetTracingHandler) GetTweetTracing(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.Tracer.Start(r.Context(), "TweetTracingHandler.GetTweetTracing")
	defer span.End()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tweetTracing, err := h.Repo.GetTweetTracing(ctx, int64(id))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(tweetTracing)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	span.SetStatus(codes.Ok, "")
	w.WriteHeader(http.StatusOK)
}

func ExtractTraceInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
