# Kewpie 队列库

[英文/English](Readme.md)

Kewpie 是一个 Go 包，实现了一个通用、动态且高效的先进先出（FIFO）队列数据结构，使用环形缓冲区。这种实现允许以摊销的常数时间复杂度进行入队和出队操作，适用于各种应用场景，包括广度优先搜索算法、事件处理系统等。

## 特点
* 通用实现：适用于任何数据类型。
* 动态调整大小：类似 Go 的 map 或切片，队列数据基于当前队列大小自动调整大小，确保高效使用内存。
* 环形缓冲区：最小化入队和出队操作的开销，因为它只是增加一个指针。
* 错误处理：为无法完成的操作提供清晰的错误信息。出错的情况很少，主要可能是因为系统内存不足。

![kewpie](https://github.com/mips171/kewpie/assets/18670565/c48e43a5-927a-4dea-82d8-85589989ff37)

*寿司排队中*

# 安装 
要在你的 Go 项目中使用 Kewpie 队列，请通过运行以下命令来安装：

```shell
go get -u github.com/mips171/kewpie
```

## 使用方法

这里有一个简单的示例，演示如何使用 Kewpie 队列来排队一些整数。

```go
package main

import (
    "fmt"
    "github.com/mips171/kewpie"
)

func main() {
    // 创建一个初始容量为10的整数队列。
    // 注意：这个 10 是可选的。如果未指定，kewpie 默认大小为 1。
    queue := kewpie.NewQueue[int](10)

    // 入队一些元素。
    queue.Enqueue(1)
    queue.Enqueue(2)
    queue.Enqueue(3)

    // 查看队首元素。
    front, err := queue.Peek()
    if err != nil {
        fmt.Println("查看队列时出错：", err)
        return
    }
    fmt.Println("队首元素是：", front)

    // 出队并打印元素。
    for {
        element, err := queue.Dequeue()
        if err != nil {
            // 遇到错误意味着队列为空。
            fmt.Println("队列为空：", err)
            break
        }
        fmt.Println("出队：", element)
    }
}
```
输出：

```sh
队首元素是：1
出队：1
出队：2
出队：3
队列为空：kewpie: 队列为空
```

# 在 BFS 中使用 Kewpie

现在让我们使用 Kewpie 为广度优先搜索排队一个自定义的 TreeNode 类型：

```go
package main

import (
    "fmt"
    "github.com/mips171/kewpie"
)

type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}


func bfs(root *TreeNode) {
    if root == nil {
        return
    }
    
    queue := kewpie.NewQueue[*TreeNode]() // 为 TreeNode 指针初始化一个新队列
    queue.Enqueue(root) // 从根节点开始 BFS
    
    while queue.Size() > 0 {
        node, err := queue.Dequeue() // 从队列中移除下一个节点
        if err != nil {
            fmt.Println("出队时出错：", err)
            return
        }
        
        fmt.Println(node.Val) // 处理当前节点
        
        // 将子节点添加到队列中以便处理
        if node.Left != nil {
            queue.Enqueue(node.Left)
        }
        if node.Right != nil {
            queue.Enqueue(node.Right)
        }
    }
}

func main() {
    // 示例二叉树：
    //     1
    //    / \
    //   2   3
    //  / \   \
    // 4   5   6
    root := &TreeNode{1,
        &TreeNode{2,
            &TreeNode{4, nil, nil},
            &TreeNode{5, nil, nil},
        },
        &TreeNode{3,
            nil,
            &TreeNode{6, nil, nil},
        },
    }
    
    fmt.Println("树的 BFS 遍历：")
    bfs(root) // [1 2 3 4 5 6]
}
```

## 


## 基准测试

```sh
go test -bench=. github.com/mips171/kewpie  -benchmem -benchtime=100x
```
此输出表明，当使用大批量大小（介于 10,000 和 100,000 之间）时，使用具有适当批量大小的函数的 Batch 版本（EnqueueBatch/DequeueBatch）会更有效，以减少 Go 所需的分配数量。

输出：

```sh
goos: darwin
goarch: arm64
pkg: github.com/mips171/kewpie
BenchmarkDequeue1-10                 100               103.0 ns/op             0 B/op          0 allocs/op
BenchmarkDequeue10-10                100               983.3 ns/op           128 B/op          4 allocs/op
BenchmarkDequeue100-10               100              1959 ns/op            1024 B/op          7 allocs/op
BenchmarkDequeue1000-10              100             11058 ns/op            8192 B/op         10 allocs/op
BenchmarkDequeue10000-10             100             88383 ns/op          131074 B/op         14 allocs/op
BenchmarkEnqueue1-10                 100               139.6 ns/op            20 B/op          0 allocs/op
BenchmarkEnqueue10-10                100               177.1 ns/op           163 B/op          0 allocs/op
BenchmarkEnqueue100-10               100               915.9 ns/op          2621 B/op          0 allocs/op
BenchmarkEnqueue1000-10              100              8073 ns/op           20971 B/op          0 allocs/op
BenchmarkEnqueue10000-10             100             75846 ns/op          167773 B/op          0 allocs/op
BenchmarkEnqueueMessages/1-10        100               912.9 ns/op           258 B/op          2 allocs/op
BenchmarkEnqueueMessages/10-10       100              7733 ns/op            2310 B/op         20 allocs/op
BenchmarkEnqueueMessages/100-10                      100             67778 ns/op           28576 B/op        200 allocs/op
BenchmarkEnqueueMessages/1000-10                     100            662272 ns/op          253876 B/op       2000 allocs/op
BenchmarkEnqueueMessages/10000-10                    100           6698946 ns/op         2286749 B/op      20000 allocs/op
BenchmarkDequeueMessages/1-10                        100               806.2 ns/op           126 B/op          1 allocs/op
BenchmarkDequeueMessages/10-10                       100              7211 ns/op            3539 B/op         27 allocs/op
BenchmarkDequeueMessages/100-10                      100             66341 ns/op           32237 B/op        211 allocs/op
BenchmarkDequeueMessages/1000-10                     100            651548 ns/op          282732 B/op       1999 allocs/op
BenchmarkDequeueMessages/10000-10                    100           6671030 ns/op         3620320 B/op      19828 allocs/op
BenchmarkEnqueueBatch/BatchSize10-10                 100               115.8 ns/op          1024 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize100-10                100               897.1 ns/op         15710 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize1000-10               100             31023 ns/op          125337 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize10000-10              100            292337 ns/op          998803 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize100000-10                     100           3284543 ns/op        16043213 B/op          0 allocs/op
BenchmarkEnqueueBatch/BatchSize1000000-10                    100          26368697 ns/op        128345702 B/op         0 allocs/op
BenchmarkDequeueBatch/BatchSize10-10                         100          93716526 ns/op        1285310464 B/op  1000020 allocs/op
BenchmarkDequeueBatch/BatchSize100-10                        100          70284867 ns/op        1291704576 B/op   100017 allocs/op
BenchmarkDequeueBatch/BatchSize1000-10                       100          66622307 ns/op        1296777216 B/op    10014 allocs/op
BenchmarkDequeueBatch/BatchSize10000-10                      100          61337613 ns/op        1287847936 B/op     1010 allocs/op
BenchmarkDequeueBatch/BatchSize100000-10                     100          60451120 ns/op        1279066112 B/op      107 allocs/op
BenchmarkDequeueBatch/BatchSize1000000-10                    100          57964251 ns/op        1235025920 B/op       14 allocs/op
PASS
ok      github.com/mips171/kewpie       114.039s
```

# 贡献
欢迎贡献！请通过 GitHub 提交问题和拉取请求，并确保你的代码遵循 Go 编码标准。

# 许可证
Kewpie 根据 MIT 许可证发布。有关更多细节，请查看 LICENSE 文件。