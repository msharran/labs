package main

import (
	"fmt"
)

type queue struct {
	array []*node
}

func (q *queue) Add(item *node) {
	q.array = append(q.array, item)
}

func (q *queue) Get() (*node, bool) {
	if len(q.array) > 0 {
		first := q.array[0]
		q.array = q.array[1:]
		return first, true
	}
	return nil, false
}

type node struct {
	value int
	left  *node
	right *node
	level int
}

//             10
//             / \
//           20   30
//            \    / \
//             32  40 45

// [
// [10]
// [20 30]
// [32 40 45]
//]

func main() {
	n10 := &node{value: 10, level: 1}
	n20 := &node{value: 20, level: 2}
	n30 := &node{value: 30, level: 2}
	n32 := &node{value: 32, level: 3}
	n40 := &node{value: 40, level: 3}
	n45 := &node{value: 45, level: 3}
	n10.left = n20
	n10.right = n30
	n20.right = n32
	n30.left = n40
	n30.right = n45

	var result [][]int
	var innerList []int
	var previousLevel int

	q := &queue{}
	q.Add(n10)
	previousLevel = 1
	for {
		item, ok := q.Get()
		if !ok {
			result = append(result, innerList)
			fmt.Println("No items in queue. exiting")
			break
		}

		fmt.Println(item.value)

		if item.level != previousLevel {
			result = append(result, innerList)
			innerList = nil
		}

		innerList = append(innerList, item.value)
		previousLevel = item.level

		if item.left != nil {
			q.Add(item.left)
		}
		if item.right != nil {
			q.Add(item.right)
		}
	}

	fmt.Println(result)
}
