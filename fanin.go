package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func toInt(done <-chan interface{}, inStream <-chan interface{}) <-chan int {
	outStream := make(chan int)
	go func() {
		defer close(outStream)
		for v := range inStream {
			select {
			case <-done:
				return
			case outStream <- v.(int):
			}
		}
	}()
	return outStream
}

func take(done <-chan interface{}, inStream <-chan interface{}, num int) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)
		for i := 0; i < num; {
			select {
			case <-done:
				return
			case outStream <- <-inStream:
				i++
			}
		}
	}()
	return outStream
}

func primeFinder(done <-chan interface{}, inStream <-chan int) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)

		for {
		Finder:
			select {
			case <-done:
				return
			case v := <-inStream:
				for i := 2; i < v; i++ {
					if v%i == 0 {
						break Finder
					}
				}
				outStream <- v
			}
		}
	}()
	return outStream
}

func do() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	rand := func() interface{} { return rand.Intn(50000000) }
	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)
	fmt.Println("Primes:")
	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}
	fmt.Printf("Search took: %v", time.Since(start))
}

func main() {
	do()
}
