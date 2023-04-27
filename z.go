// Package z allows you to build simple CLI tools that process STDIN lines
// concurrently. Implement a factory and a task. Then Run() your factory.
package z

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

var Grs = 100 // number of concurrent goroutines processing the tasks

type Task interface {
	Process()
	Print()
}

type Factory interface {
	Make(line string) Task
}

func Run(f Factory) {
	var wg sync.WaitGroup

	in := make(chan Task)

	wg.Add(1)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			in <- f.Make(s.Text())
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "z: reading STDIN: %s", err)
		}
		close(in)
		wg.Done()
	}()

	out := make(chan Task)

	for i := 0; i < Grs; i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				t.Process()
				out <- t
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	for t := range out {
		t.Print()
	}
}
