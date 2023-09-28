package workers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WorkerResult struct {
	Line    string
	LineNum int
	Path    string
}

type AllResults struct {
	FinalRes []WorkerResult
}

func NewWorkerResult(line string, lineNum int, path string) WorkerResult {
	return WorkerResult{
		Line:    line,
		LineNum: lineNum,
		Path:    path,
	}
}

func FindSearchTerm(filePath string, searchTerm string) *AllResults {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("%s was not found", filePath)
		return nil
	}

	workerResult := AllResults{make([]WorkerResult, 0)}

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchTerm) {
			r := NewWorkerResult(scanner.Text(), lineNumber, filePath)
			workerResult.FinalRes = append(workerResult.FinalRes, r)
		}
		lineNumber++

	}

	if len(workerResult.FinalRes) == 0 {
		return nil
	}
	return &workerResult
}
