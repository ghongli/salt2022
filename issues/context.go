package issues

import (
	"context"
	"fmt"
	"time"
)

func ctxTimeoutDone(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("handler", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}
}