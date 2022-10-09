# pkg说明

### algs

主要提供优先级队列的功能。通过```algs.GetPQ()```创建队列。目前队列元素只支持string类型，使用方式如下所示，priority值越小，优先级越高：

```go
pq := algs.GetPQ()
	item2 := algs.GetItem("test2", 2)
	item1 := algs.GetItem("test1", 1)
	item3 := algs.GetItem("test3", 3)
	item0 := algs.GetItem("test0", 0)
	pq.Push(item2)
	pq.Push(item1)
	pq.Push(item3)
	pq.Push(item0)
	for pq.Len() > 0 {
		item, ok := pq.PopItem()
		fmt.Printf("item: %#v,ok: %v\n", item, ok)
	}
```

```bash
item: &algs.Item{value:"test0", priority:0, index:-1},ok: true
item: &algs.Item{value:"test1", priority:1, index:-1},ok: true
item: &algs.Item{value:"test2", priority:2, index:-1},ok: true
item: &algs.Item{value:"test3", priority:3, index:-1},ok: true
```
