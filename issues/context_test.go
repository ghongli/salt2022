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
	
	t.Run("timeout duration gt parent ctx", func(t *testing.T) {
		go ctxTimeoutDone(ctx, 1500 * time.Millisecond)
		
		select {
		case <-ctx.Done():
			fmt.Println("main", ctx.Err())
		}
	})
	
	t.Run("timeout duration lt parent ctx", func(t *testing.T) {
		go ctxTimeoutDone(ctx, 500 * time.Millisecond)
		
		select {
		case <-ctx.Done():
			fmt.Println("main", ctx.Err())
		}
	})
	
}