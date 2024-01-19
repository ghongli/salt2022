package main

import (
	"fmt"
	
	"golang.org/x/net/context"
)

func main() {
	type favContextKey string
	
	f := func(ctx context.Context, k any) {
		if v := ctx.Value(k); v != nil {
			fmt.Printf("%s found value: %v \n", k, v)
			return
		}
		fmt.Printf("%s not found value \n", k)
	}
	
	// 避免多层 valueCtx 查找的时候冲突，得到不是想要的结果。防止不同包的 key 冲突，同一个包防不了的。
	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go2")
	ctx = context.WithValue(ctx, "language", "Go")
	
	f(ctx, k)
	f(ctx, "language")
	f(ctx, "color")
}
