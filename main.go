package main

import (
	"flag"
	"os"
	"text/template"
)

type TypeInfo struct {
	Package string
	Key string
	Value string
}

func main() {
	var t TypeInfo
	flag.StringVar(&t.Package, "package", "", "name of package")
	flag.StringVar(&t.Key, "key", "", "key to use for LRU Cache")
	flag.StringVar(&t.Value, "value", "", "value held in LRU Cache")
	flag.Parse()

	temp := template.Must(template.New("lrucache").Parse(cachTemplate))
	temp.Execute(os.Stdout, t)
}

var cachTemplate = `
package {{.Package}}

import "fmt"

type MissingKey struct {
	key {{.Key}}
}

func (m MissingKey) Error() string {
	return fmt.Sprintf("no key %v found", m.key)
}

type node struct {
	Key   {{.Key}}
	Value {{.Value}}
	Next  *node
	Prev  *node
}

func newNode(key {{.Key}}, value {{.Value}}) *node {
	return &node{Key: key, Value: value}
}

func (n *node) String() string {
	return fmt.Sprintf("<Node: Key:%v, Value:%v>", n.Key, n.Value)
}

type linkedList struct {
	Head *node
	Tail *node
}

func newLinkedList() *linkedList {
	return &linkedList{nil, nil}
}

func (l *linkedList) Add(n *node) {
	if l.Head == nil {
		l.Head = n
		l.Tail = n
	} else {
		head := n
		head.Next = l.Head
		l.Head.Prev = head
		l.Head = head
	}
}

func (l *linkedList) Remove(n *node) {
	if n == l.Head {
		if n.Next != nil {
			l.Head = n.Next
		} else {
			l.Head = nil
		}
		return
	}
	if n == l.Tail {
		if n.Prev != nil {
			l.Tail = n.Prev
		}
	}
	nextNode := n.Next
	prevNode := n.Prev
	if nextNode != nil {
		nextNode.Prev = prevNode
	}
	if prevNode != nil {
		prevNode.Next = nextNode
	}
}

type LRUCache struct {
	Capacity int
	ll       *linkedList
	nodeMap  map[{{.Key}}]*node
}

func NewCache(cap int) *LRUCache {
	return &LRUCache{
		Capacity: cap,
		ll:       newLinkedList(),
		nodeMap:  make(map[string]*node),
	}
}

func (c *LRUCache) Get(key {{.Key}}) ({{.Value}}, error) {
	node, ok := c.nodeMap[key]
	if ok {
		result := node.Value
		c.ll.Remove(node)
		c.ll.Add(node)
		return result, nil
	}
	return nil, MissingKey{key: key}
}

func (c *LRUCache) Put(key {{.Key}}, value {{.Value}}) {
	node, ok := c.nodeMap[key]
	if ok {
		node.Value = value
		c.ll.Remove(node)
		c.ll.Add(node)
	} else {
		if len(c.nodeMap) == c.Capacity {
			tail := c.ll.Tail
			delete(c.nodeMap, tail.Key)
			c.ll.Remove(tail)
		}
		node := newNode(key, value)
		c.nodeMap[key] = node
		c.ll.Add(node)
	}
}

`
