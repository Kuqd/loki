package queryrange

import (
	"context"
)

// StepAlignMiddleware aligns the start and end of request to the step to
// improved the cacheability of the query results.
var StepAlignMiddleware = MiddlewareFunc(func(next Handler) Handler {
	return stepAlign{
		next: next,
	}
})

type stepAlign struct {
	next Handler
}

func (s stepAlign) Do(ctx context.Context, r Request) (Response, error) {
	return s.next.Do(ctx, r.WithStartEnd((r.GetStart()/r.GetStep())*r.GetStep(), (r.GetEnd()/r.GetStep())*r.GetStep()))
}
