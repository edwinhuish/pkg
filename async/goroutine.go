package async

import (
	"context"
	"fmt"
	"log"
	"time"
)

// 用于在新的goroutine中执行给定函数并捕获可能的panic的函数。
// 如果提供了错误处理函数，它将被用于处理捕获的panic错误。
// 参数:
//
//	task: 要执行的函数，没有参数和返回值。
//	finally: 可选的最终处理函数，error 为 nil 代表正常结束。
func DoWithContext(ctx context.Context, task func(), finally ...func(err error)) {

	go func() {

		var taskErr error = nil

		// 根据是否提供了错误处理函数来决定使用哪一个。
		if len(finally) == 0 {
			// 默认的错误处理函数，将错误记录到日志。
			finally = append(finally, func(err error) {
				if err != nil {
					log.Printf("async error: %v\n", err)
				}
			})
		}

		defer func() {
			for _, f := range finally {
				f(taskErr)
			}
		}()

		defer func() {
			// 检查是否有panic发生，如果有，则根据情况调用错误处理函数。
			if err := recover(); err != nil {

				// 根据recover返回的错误类型，调用错误处理函数。
				switch v := err.(type) {
				case error:
					taskErr = v
				default:
					// 如果不是error类型，创建一个error类型
					taskErr = (fmt.Errorf("%+v", v))
				}
			}
		}()

		// 如果提供了上下文，则使用上下文来等待任务完成。
		select {
		case <-ctx.Done():
			return
		default:
			// 执行传入的函数。
			task()
		}
	}()
}

func DoWithTimeout(timeout time.Duration, task func(), finally ...func(err error)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	finally = append(finally, func(err error) {
		cancel()
	})

	DoWithContext(ctx, task, finally...)
}

func Do(task func(), finally ...func(err error)) {
	DoWithTimeout(defaultTimeout, task, finally...)
}
