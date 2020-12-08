学习笔记

1. 基础库中的errors.New返回了errorString对象的指针，主要是为了避免在error判等时造成不同的异常相等的情况
    ```go
    fmt.Println(errors.New("123") == errors.New("123")) // false
    fmt.Println(New("123") == New("123")) //true
    
    type errString struct {
        s string
    }
    
    func New(msg string) errString {
        return errString{s: msg}
    }
    
    func (e *errString) Error() string {
        return e.s
    }
    ```
2. Java异常的处理，应该尽量根据调用方法所抛出的异常做出正确的处理，目前自己滥用严重
    
    golang中使用多返回值的形式，将正确的返回值和error分离开。两者一般按互斥处理。
    
    一般而言，只在一个地方处理error
    
3. panic
    - 在web应用主线程中，一般使用中间件recover处理panic。
    - 在协程中，一般将go关键字包装为一个方法，然后recover处理panic
        ```go
        func Go(run func()) {
            go func() {
                defer func() {
                    if err := recover(); err != nil {
                        errorHandler(err)
                    }
                }()
                run()
            }()
        }
          
        ```

4. Sentinel Error 哨兵错误
    
    内部库对外部直接提供公有的Error，由外部直接调用。但是这种会增强调用方对内部库的依赖。

5. 断言
    - 断言类型
        
        我们可以根据上下文自定义自己的Error
        
        ```go
        switch err := err.(type) {
        case nil:
        case *MyError:
        	fmt.Println("xxxx")
        default:
        }
        ```
    - 断言行为
    
        Opaque Error，将Sentinel Error封装在内部库中，不对外暴露，只提供公有方法进行错误判断
        
        ```go
        type temporary interface {
      	    Temporary() bool
        }
        
        func IsTemporary(err error) bool {
            te, ok := err.(temporary)
            return ok && te.Temporary()
        }
        ```

6. 优化error

    bufio.NewReader(); br.ReadString("\n"); bufio内部使用了NewScanner优化了代码处理
    
    ```go
    type errWriter struct {
    	io.Writer
    	err error
    }
    
    // 注: 执行完所有操作之后，直接判断e.err即可
    func (e *errWriter) Write(buf []byte) (int, error) {
    	// 如果已经有错误，就直接返回错误
    	if e.err != nil {
    		return 0, e.err
    	}
    	var n int
    	// 把第一次的错误存起来
    	n, e.err = e.Writer.Write(buf)
    	return n, nil
    }
    ```

7. wrap errors

    github.com/pkg/errors开源包封装了错误堆栈，可以使用%+v打印堆栈信息
    
    我们一般会在如下地方包装error
    - 第三方库error
    - 基础库error
    - 封装的kit库error
    
8. Go 1.13

    - Unwrap
        
        返回原始包装的error。
    - Is
    
        递归断言Error类型。
    - As
        
        递归取出具体的错误类型，封装到变量中。
    