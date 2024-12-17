package async

import (
	"fmt"
	"log"
)

// 用于在新的goroutine中执行给定函数并捕获可能的panic的函数。
// 如果提供了错误处理函数，它将被用于处理捕获的panic错误。
// 参数:
//
//	task: 要执行的函数，没有参数和返回值。
//	finally: 可选的最终处理函数，error 为 nil 代表正常结束。
func Do(task func(), finally ...func(err error)) {

	go func() {

		var taskErr error = nil

		var _finally func(err error)
		// 根据是否提供了错误处理函数来决定使用哪一个。
		if len(finally) > 0 {
			_finally = finally[0]
		} else {
			// 默认的错误处理函数，将错误记录到日志。
			_finally = func(err error) {
				if err != nil {
					log.Printf("async error: %v\n", err)
				}
			}
		}

		defer _finally(taskErr)

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

		// 执行传入的函数。
		task()
	}()
}
