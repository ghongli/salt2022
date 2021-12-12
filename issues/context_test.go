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
		
		if <-ctx.Done(); true {
			fmt.Println("main", ctx.Err())
		}
		
		// output:
		// handler context deadline exceeded
		// main context deadline exceeded
	})
	
	t.Run("timeout duration lt parent ctx", func(t *testing.T) {
		go ctxTimeoutDone(ctx, 500 * time.Millisecond)
		
		if <-ctx.Done(); true {
			fmt.Println("main", ctx.Err())
		}
		
		// output:
		// process request with 500ms
		// main context deadline exceeded
	})
	
}