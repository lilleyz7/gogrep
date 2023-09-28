package main

import (
	"cgrep/workers"
	"cgrep/worklist"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
)

func discoverDirectories(wl *worklist.Worklist, path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			nextPath := filepath.Join(path, entry.Name())
			discoverDirectories(wl, nextPath)
		} else {
			wl.AddJob(worklist.NewJob(filepath.Join(path, entry.Name())))
		}
	}
}

var args struct {
	SearchTerm      string `arg:"positional,required"`
	SearchDirectory string `arg:"positional"`
}

func main() {
	arg.MustParse(&args)
	fmt.Println(args.SearchTerm)
	fmt.Println(args.SearchTerm)

	var workersGroup sync.WaitGroup

	wl := worklist.NewWorkList(100)
	results := make(chan workers.WorkerResult, 100)
	totalWorkers := 10

	workersGroup.Add(1)

	go func() {
		defer workersGroup.Done()
		discoverDirectories(&wl, args.SearchDirectory)
		wl.CompleteWork(totalWorkers)
	}()

	for i := 0; i < totalWorkers; i++ {
		workersGroup.Add(1)
		go func() {
			defer workersGroup.Done()
			for {
				workEntry := wl.NextJob()
				workerPath := workEntry.GetPath()
				if workerPath != "" {
					workResult := workers.FindSearchTerm(workerPath, args.SearchTerm)
					if workResult != nil {
						for _, r := range workResult.FinalRes {
							results <- r
						}
					}
				} else {
					return
				}
			}
		}()
	}

	blockWorkersWg := make(chan struct{})
	go func() {
		workersGroup.Wait()
		close(blockWorkersWg)
	}()

	var displayWg sync.WaitGroup
	displayWg.Add(1)

	go func() {
		for {
			select {
			case r := <-results:
				fmt.Printf("%v [%v]: %v \n", r.Path, r.LineNum, r.Line)
			case <-blockWorkersWg:
				if len(results) == 0 {
					displayWg.Done()
					return
				}
			}
		}

	}()
	displayWg.Wait()

}
