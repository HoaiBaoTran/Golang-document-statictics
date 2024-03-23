package main

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hoaibao/statistics-document/statistics"
)

type Result struct {
	lineCount, wordCount, charCount, totalWordLength int
	frequency                                        map[string]int
}

var (
	result Result = Result{
		frequency: make(map[string]int),
	}
	chunkNumber = 5
)

func merge(workerChannels ...<-chan Result) <-chan Result {
	var wg sync.WaitGroup
	out := make(chan Result)

	copyToOutput := func(c <-chan Result) {
		for item := range c {
			out <- item
		}
		wg.Done()
	}

	wg.Add(len(workerChannels))
	for _, c := range workerChannels {
		go copyToOutput(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func handleChunk(chunk []string) <-chan Result {
	chunkResultChan := make(chan Result)
	go func() {
		frequencyChan := statistics.CountFrequencyFromLine(strings.Join(chunk, " "))
		defer close(chunkResultChan)
		chunkResult := Result{}
		for _, line := range chunk {
			wordCount, charCount := statistics.WordAndCharCount(line)
			chunkResult.wordCount += wordCount
			chunkResult.charCount += charCount
		}
		chunkResult.frequency = <-frequencyChan
		chunkResultChan <- chunkResult
	}()
	return chunkResultChan
}

func splitFile(lines []string, chunkNumber int, wg *sync.WaitGroup) [][]string {
	chunks := make([][]string, chunkNumber)
	lengthPerChunks := len(lines) / chunkNumber

	for i := 0; i < chunkNumber; i++ {
		go func(i int) {
			var chunkLines []string
			if i == chunkNumber-1 {
				chunkLines = lines[i*lengthPerChunks:]
			} else {
				chunkLines = lines[i*lengthPerChunks : (i+1)*lengthPerChunks]
			}
			chunks[i] = chunkLines
			wg.Done()
		}(i)
	}

	return chunks
}

func readingFile(filePath string) ([]string, int) {
	file, err := os.Open(filePath)
	statistics.CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
		lineCount++
	}
	statistics.CheckError(scanner.Err())
	return lines, lineCount
}

func splitFileMethod(filePath string, startTime time.Time) {
	var wg sync.WaitGroup

	lines, lineCount := readingFile(filePath)

	result.lineCount = lineCount

	wg.Add(chunkNumber)
	chunks := splitFile(lines, chunkNumber, &wg)
	wg.Wait()

	results := make([]<-chan Result, len(chunks))
	for index, chunk := range chunks {
		results[index] = handleChunk(chunk)
	}

	returnChan := merge(results...)
	for i := range returnChan {
		result.charCount += i.charCount
		result.wordCount += i.wordCount
		for key, value := range i.frequency {
			if _, isKey := result.frequency[key]; !isKey {
				result.totalWordLength += len(key)
				result.frequency[key] = value
			} else {
				result.frequency[key] += value
			}
		}
	}

	executionTime := time.Since(startTime)

	statistics.WriteResultToFile(
		result.lineCount,
		result.wordCount,
		result.charCount,
		float64(result.totalWordLength)/float64(len(result.frequency)),
		result.frequency,
		executionTime,
	)

}

func main() {
	startTime := time.Now()

	// filePath := "../documents/requirement_en.txt"
	// filePath := "../documents/file.txt"
	filePath := "../documents/file_1MB.txt"
	splitFileMethod(filePath, startTime)
}
