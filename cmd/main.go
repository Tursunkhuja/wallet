package main

import (
	"log"
	"sync"
)

func main() {
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}
	parts := 10
	size := len(data) / parts
	channels := make([]<-chan int, parts)

	for i := 0; i < parts; i++ {
		ch := make(chan int)
		channels[i] = ch
		go func(ch chan<- int, data []int) {
			defer close(ch)
			sum := 0
			for _, v := range data {
				sum += v
			}
			ch <- sum
		}(ch, data[i*size:(i+1)*size])

	}

	total := 0
	for value := range merge(channels) {
		total += value
	}

	log.Print(total)
}

func merge(channels []<-chan int) <-chan int {
	wg := sync.WaitGroup{}
	wg.Add(len(channels))
	merged := make(chan int)

	for _, ch := range channels {

		go func(ch <-chan int) {
			defer wg.Done()
			for val := range ch {
				merged <- val
			}
		}(ch)
	}

	go func() {
		defer close(merged)
		wg.Wait()
	}()

	return merged
}
