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

## Usage

Here's a simple example demonstrating how to use the Kewpie queue to queeue some integers.

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
Output:
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

## 基准测试

```sh
go test -bench=. github.com/mips171/kewpie -benchmem
```

输出：

```sh
goos: darwin
goarch: arm64
pkg: github.com/mips171/kewpie
BenchmarkEnqueue-10     176774624                6.216 ns/op          24 B/op          0 allocs/op
BenchmarkDequeue-10     68080922                17.43 ns/op           15 B/op          0 allocs/op
PASS
ok      github.com/mips171/kewpie       20.812s
```

# 贡献
欢迎贡献！请通过 GitHub 提交问题和拉取请求，并确保你的代码遵循 Go 编码标准。

# 许可证
Kewpie 根据 MIT 许可证发布。有关更多细节，请查看 LICENSE 文件。