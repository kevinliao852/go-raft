package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	ch     chan int
	id     int
	status string
}

func node(node *Node, nodes *[]Node, mu *sync.Mutex, wg *sync.WaitGroup) {

	mu.Lock()
	if node.id == 0 {
		node.status = "leader"
	} else {
		node.status = "follower"
	}
	wg.Done()
	mu.Unlock()

	cnt := rand.Float64()*1000 + float64(node.id)

	// if node.id == 0 {
	// 	cnt = float64(8000)
	// }
	// timeout := time.After(time.Duration(cnt))

	for {
		fmt.Println(node.id, node, (*nodes)[node.id])
		select {
		// case <-timeout:
		// 	mu.Lock()
		// 	node.status = "leader"
		// 	mu.Unlock()
		case data := <-node.ch:
			if node.status == "follower" {
				//fmt.Println("Node", node.id, "received data from Node", data)
				// Response to Leader
				(*nodes)[data].ch <- node.id
			} else {
				// leader receives data
				//fmt.Println("Leader Node", node.id, "received data from Node", data)
			}
			// todo: need to response
		default:
			if node.status == "leader" {
				go func() {
					for i := range *nodes {
						if i != node.id {
							//fmt.Println("Node", node.id, "sends heartbeat to Node", i)
							(*nodes)[i].ch <- node.id
						}
					}
				}()
			}
		}

		//fmt.Println(node.id, node.status, cnt)

		time.Sleep(time.Millisecond * time.Duration(cnt))
	}
}

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	nodes := []Node{}
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		n := Node{make(chan int), i, "follower"}
		wg.Add(1)
		nodes = append(nodes, n)
		// fmt.Println("xxx", nodes)

		go node(&nodes[i], &nodes, &mu, &wg)
		wg.Wait()
	}

	for j := range ch {
		fmt.Println(j)

	}
}
