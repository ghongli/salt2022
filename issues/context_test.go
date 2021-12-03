package issues

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTimeoutDone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(),1*time.Second)
	defer cancel()
	
	handler := func(ctx context.Context, duration time.Duration) {
		select {
		case <-ctx.Done():
			fmt.Println("handler", ctx.Err())
		case <-time.After(duration):
			fmt.Println("process request with", duration)
		}
	}
	
	go handler(ctx, 500 * time.Millisecond)
	
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}
