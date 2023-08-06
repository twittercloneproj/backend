package data

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type TweetTracingRepository struct {
	Tracer        trace.Tracer
	TweetsTracing map[int64]TweetTracing
}

func (r TweetTracingRepository) GetTweetTracing(ctx context.Context, id int64) (TweetTracing, error) {
	ctx, span := r.Tracer.Start(ctx, "TweetTracingRepository.GetTweetTracing")
	defer span.End()

	tweetTracing, ok := r.TweetsTracing[id]
	if !ok {
		err := errors.New(fmt.Sprintf("tweetTracing ID = %d not found", id))
		span.SetStatus(codes.Error, err.Error())
		return TweetTracing{}, err
	}
	return tweetTracing, nil
}
