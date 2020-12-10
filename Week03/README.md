学习笔记

1. CSP模型(communicating sequential processes) 顺序通信过程

    多个goroutine可以通过管道(channel)传输消息。
    
    在Java中，通过共享(资源)内存的形式进行线程间通信，在GO中提供了另外一种通信方式，就是CSP

2. 并发模型

    GO调度器不是抢占式调度器，而是协作式调度器。

    [GO并发模型](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)
    
    ![image.png](https://i.loli.net/2020/12/06/tQ3dFNZzMsrJbqK.png)
    
    - 逻辑处理器(P)，每个虚拟内核提供一个逻辑处理器
    - OS Thread(M)，每一个P被分配一个M，OS将M放置在核心中处理逻辑
    - goroutines(G)，OS线程和内核的关系 类似于 goroutine和OS线程
    - GRQ(全局运行队列)
    - LRQ(本地运行队列)，每一个P分配一个LRQ，管理那些分配给这个P执行的G，M会对这些G依次进行上下文切换
    
    M与G是一个N:M的模型
    
    1. 异步网络系统调用
    
        ![image.png](https://i.loli.net/2020/12/06/seLaPxmJKNBtGod.png)
        
        当M上绑定的G1需要进行异步网络系统调度时，该G1将会移动到网络轮询器，然后LRQ上的G2会进行上下文切换，绑定到M上处理逻辑.
        
        防止G1进行网络系统调度时会阻塞M
        
        ![image.png](https://i.loli.net/2020/12/06/i8lN3f5OMaCKPU1.png)
        
        当G1完成系统调度时，G1会重新放入LRQ队尾
    
    2. 同步系统调用
       
        ![image.png](https://i.loli.net/2020/12/06/tdH8oTjN1VwSs2p.png)
        
        G在执行同步系统调用时，会阻塞M，此时G1将M1阻塞，在M1上绑定的P将会重新绑定到新的M2，然后有M2进行调度LRQ中的G
        G1执行完之后会重新放入LRQ队尾。
    
    3. 工作窃取(working-stealing)
       
        ```go
        runtime.schedule() {
           // only 1/61 of the time, check the global runnable queue for a G.
           // if not found, check the local queue.
           // if not found,
           //     try to steal from other Ps.
           //     if not, check the global runnable queue.
           //     if not found, poll network.
        }
        ```

3. 方法设计时，要注意函数执行的时间是否会过长；是否需要使用异步的形式获取结果；如果异步过程中符合条件的结果已经出现，如何终止方法执行
    
    ```go
    // 全量获取目录列表，如果目录过多，耗费时间会很长
    func ListDirectory(dir string) ([]string, error)
    
    // 通过channel异步获取目录，并放入channel中
    func ListDirectory(dir string) chan string
   
    // 可以通过一个方法回调判断是否满足跳出的条件
    func ListDirectory(dir string, fn func(string) bool) chan string
    ```
    
    filepath.WalkDir 也是类似的模型，如果函数启动 goroutine，则必须向调用方提供显式停止该goroutine 的方法。通常，将异步执行函数的决定权交给该函数的调用方通常更容易。
    
4. 使用goroutine时我们必须考虑两个问题
    - When will it terminate? 什么时候终止
    - What could make it terminate? 什么能让它终止
    
    即控制goroutine的整个生命周期。
    
    1. 使用channel控制goroutine的关闭
        ```go
        // 使用stop、done两个channel控制goroutine的关闭
        func (g *Group) Run() error {
           if len(g.fns) == 0 {
               return nil
           }
           stop := make(chan struct{})
           done := make(chan error, len(g.fns))
           for _, fn := range g.fns {
               go func(fn Run) {
                   // 在fn中使用 <-stop来控制资源的释放和关闭
                   done <- fn(stop)
               }(fn)
           }
           var err error
           for i := 0; i < cap(done); i++ {
               if i == 0 {
                   // 返回第一个error
                   err = <-done
                   close(stop)
               } else {
                   <-done
               }
           }
           close(done)
           return err
        }
        ```
    2. 使用WaitGroup和channel对goroutine进行控制
        WaitGroup可以控制一组goroutine执行完毕之后，在处理剩余逻辑。
        ```go
        func PoolExecute(int n) {
           ch := make(chan bool, n)
           var g sync.WaitGroup
           for i := 0; i < n; i++ {
               g.Add(1)
               ch <- true
               // 启动协程处理业务
               go func() {
                   // 在最后使WaitGroup减1
                   defer g.Done()
                   // DoSomething
                   //使用channel控制最多只能有n个协程
                   <-ch
               }
           }
           g.Wait()
        }
        ```
    3. 使用超时对goroutine的执行时间进行控制

5. 内存模型
    
    1. Happen-Before
    
        可见性
        
6. sync包

    1. Once 双重检查锁，执行且仅执行一次f函数
        
7. context包

    在服务器请求的生命周期中的function链应该传递Context
    
    可以选择性的使用WithCancel、WithDeadline、WithTimeout、WithValue等包装上下文
    
    父Context被取消之后，所有派生的子Context都会被取消
    
    ```go
    type Context interface {
       // 返回可以被取消的上下文-任务被取消时时间 如果可以被取消，返回true
       Deadline() (deadline time.Time, ok bool)
       // 上下文被主动取消、超时时，会关闭channel
       //   如果context不能被取消，则返回nil
       Done() <-chan struct{}
       // 如果Done没有被close，返回nil
       // 如果Context被取消之后，返回Canceled
       // 如果Context超时，返回DeadlineExceeded（这里的error是一个结构体，而不是指针）
       Err() error
       // 存储请求层面的key-value
       Value(key interface{}) interface{}
    }
    ```
    
    1. 应用
    
        - 超时控制
        - 主动取消
        - 全局携带（与业务无关的）k-v值
        
        

## 未完待续


1. goroutine生命周期要清楚，避免goroutine泄露
    通过context控制超时
    通过channel发送消息
    通过close channel
    
2. 启动goroutine交给调用方

3. https://golang.org/ref/mem
