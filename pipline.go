package main

import (
	"fmt"
	"strconv"
	"time"
)

func Batch(data []int) {
	adder := func(data []int, num int) []int {
		for i, _ := range data {
			data[i] += num
		}
		return data
	}
	stringer := func(data []int) []string {
		out := make([]string, len(data))
		for i, d := range data {
			out[i] = strconv.Itoa(d)
		}
		return out
	}
	addSuffix := func(data []string, suffix string) []string {
		//out := make([]string, len(data))
		for i, _ := range data {
			data[i] = data[i] + suffix
		}
		return data
	}

	batch := addSuffix(stringer(adder(data, 4)), "suff!!!")
	s := time.Now()
	for _, _ = range batch {
		//fmt.Println(d)
	}
	fmt.Println("done", time.Now().Sub(s))
}

func Pipeline(data ...int) {
	generator := func(data ...int) <-chan int {
		stream := make(chan int)
		go func() {
			defer close(stream)
			for _, d := range data {
				stream <- d
			}
		}()
		return stream
	}
	adder := func(inStream <-chan int, num int) <-chan int {
		outStream := make(chan int)
		go func() {
			defer close(outStream)
			for d := range inStream {
				outStream <- d + num
			}
		}()
		return outStream
	}
	stringer := func(inStream <-chan int) <-chan string {
		outStream := make(chan string)
		go func() {
			defer close(outStream)
			for d := range inStream {
				outStream <- strconv.Itoa(d)
			}
		}()
		return outStream
	}
	addSuffix := func(inStream <-chan string, suffix string) <-chan string {
		outStream := make(chan string)
		go func() {
			defer close(outStream)
			for d := range inStream {
				outStream <- d + suffix
			}
		}()
		return outStream
	}
	pipeline := addSuffix(stringer(adder(generator(data...), 4)), "suff!!!")
	s := time.Now()
	for _ = range pipeline {
		//fmt.Println(d)

	}
	fmt.Println("done", time.Now().Sub(s))
}

func main() {
	data := make([]int, 10000000)
	for i, _ := range data {
		data[i] = i
	}
	Pipeline(data...)
	//Batch(data)

}
