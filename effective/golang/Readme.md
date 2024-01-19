effective golang
---

#### issues

1. go 自己内部做完 gc 后，内存还是会被进程占用着，在什么情况下，这部分内存才会真正归还给操作系统？
    gc 之后还有 scavenge，这个时候就涉及 MADV_FREE VS MADV_DONTNEED 的问题，go1.12-1.15 用 MADV_FREE，不会立即释放，会有 RSS 虚高，导致 OOM 的问题。
    
2. 为什么 context 的 key 最好不要使用内置类型？
   不使用内置类型，如果创建相同的 key，如 favContextKey("language")，只要字符串一样，都会碰撞，和直接使用 "language" 没有区别？
   
   避免多层 valueCtx 查找的时候冲突，得到不是想要的结果。防止不同包的 key 冲突，同一个包防不了的。
   ```go
   // 使用 WithValue 时，string 作为 key 导致命名冲突
   
   package main
   
   import (
      "context"
      "fmt"
   )
   
   func main() {
      ctx := context.Background()
      ctx = context.WithValue(ctx, "key", "value1")
      fmt.Println(ctx.Value("key")) // Output: value1
   
      ctx = context.WithValue(ctx, "key", "value2")
      fmt.Println(ctx.Value("key")) // Output: value2
   }
   // 由于两次使用相同的字符串 key，第二次 WithValue 会覆盖之前的值，导致出现命令冲突。
   ```
   
   ```go
   // 不使用内置类型，使用自定义类型，不会产生命名冲突，覆盖值的情况
   package main
   
   import (
      "context"
      "fmt"
   )

   func main() {
      // type customKey string
      // ctx := context.Background()
      // ctx = context.WithValue(ctx, customKey("key"), "value1")
      // fmt.Println(ctx.Value(customKey("key"))) // Output: value1
   
	  // 自定义类型，可发防止相同值不会覆盖。
      // ctx = context.WithValue(ctx, customKey("key"), "value2")
      // fmt.Println(ctx.Value(customKey("key"))) // Output: value2
   
	  // 使用自定义类型作为 context key，确保其唯一性和清晰性
	  type favContextKey string
   
	  f := func(ctx context.Context, k favContextKey) {
        if v := ctx.Value(k); v != nil {
            fmt.Printf("%s found value: %v \n", k, v)
            return
        }
        fmt.Printf("%s not found value \n", k)
	  }
   
	  k := favContextKey("language")
	  ctx := context.WithValue(context.Background(), k, "Go1.20")
	  
	  f(ctx, k)
	  f(ctx, favContextKey("color"))
   }
   ```
   
   ```go
   // 使用不同的结构体
   func NewContext(ctx context.Context, info jwt.Claims) context.Context {
	   return context.WithValue(ctx, authKey{}, info)
   }
   
   func FromContext(ctx context.Context) (token jwt.Claims, ok bool) {
	   return ctx.Value(authKey{}).(jwt.Claims)
   }
   ```

3. 

#### wiki

1. [SliceTricks](https://github.com/golang/go/wiki/SliceTricks)

---
[0]: https://github.com/golang/go/wiki "golang wiki"