package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// context 上下文环境
//  同一个请求中所有处理,包括共享数据、子goroutine资源释放
//查看pprof堆栈使用 http://127.0.0.1:8989/debug/pprof
/* output
2019/05/22 15:30:54 B: 12
2019/05/22 15:30:54 C: python
2019/05/22 15:30:57 C done
2019/05/22 15:30:57 A done
2019/05/22 15:30:57 B done
*/

//传递的自定义键类型
// 避免与作用域内置类型发生碰撞
type key string

var UseKey = key("Name")

func main() {
	go http.ListenAndServe(":8989", nil)

	// 启动一个根的上下文ctx
	ctx, cancel := context.WithCancel(context.Background())
	// 在一个ctx中共享数据,k/v
	ctx = context.WithValue(ctx, UseKey, 12)
	go func() {
		time.Sleep(3 * time.Second)
		// 当根ctx退出,子ctx收到Done信号
		cancel()
	}()

	log.Println(A(ctx))
	select {}
}

func A(ctx context.Context) string {
	go log.Println(B(ctx))
	select {
	case <-ctx.Done():
		return "A done"
	}
	return ""
}

func B(ctx context.Context) string {
	// 根ctx传递的共享数据
	log.Println("B:", ctx.Value(UseKey))
	// 重设置共享数据
	ctx = context.WithValue(ctx, UseKey, "python")
	go log.Println(C(ctx))
	select {
	case <-ctx.Done():
		return "B done"
	}
	return ""
}

func C(ctx context.Context) string {
	// B传递的共享数据
	log.Println("C:", ctx.Value(UseKey))
	select {
	case <-ctx.Done():
		return "C done"
	}
	return ""
}
