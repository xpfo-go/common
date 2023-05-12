# Stack

## 一 介绍

### 链表实现得LIFO得泛型栈结构，比切片实现得栈性能好一点

## 二 使用

``` go

    s := Stack[int]{}
    s.Push(1) // 入栈
    s.Pop() // 出栈 返回栈顶元素
    s.Size()  // 栈长度
    s.Empty() // 栈是否为空
    s.Top()  // 返回栈顶元素
    s.Bottom()  // 返回栈底元素

```