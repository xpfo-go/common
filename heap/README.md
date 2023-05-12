# Heap

## 一 介绍

### 基于切片实现得泛型堆结构

## 二 使用

``` go
    type node struct {
        a int
    }
    
    // 第一个参数是表示初始化得数据  第二个参数是指定最大堆还是最小堆 这里传nil得话会panic
    h := NewHeap[*node](nil, func(a, b *node) bool {
        return a.a < b.a
    })
	
    h.Pop()  // 出堆
    h.Size()  // 堆长
    h.Empty()  // 是否为空
    h.Push(&node{a: 12})  // 入堆  
```