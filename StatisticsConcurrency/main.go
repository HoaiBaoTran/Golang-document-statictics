package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hoaibao/statistics-document/statistics"
)

type Result struct {
	lineCount, wordCount, charCount, totalWordLength int32
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
		defer close(chunkResultChan)
		chunkResult := Result{}
		for _, line := range chunk {
			wordCount, charCount := statistics.WordAndCharCount(line)
			chunkResult.wordCount += int32(wordCount)
			chunkResult.charCount += int32(charCount)
		}
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

func readingFile(filePath string) []string {
	file, err := os.Open(filePath)
	statistics.CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		result.lineCount++
	}
	statistics.CheckError(scanner.Err())

	return lines
}

func splitFileMethod(filePath string, startTime time.Time) {
	var wg sync.WaitGroup

	lines := readingFile(filePath)

	wg.Add(chunkNumber + 1)
	go func() {
		result.frequency, result.totalWordLength = statistics.CountFrequencyFromLine(lines, &wg)
	}()

	chunks := splitFile(lines, chunkNumber, &wg)

	results := make([]<-chan Result, len(chunks))
	for index, chunk := range chunks {
		results[index] = handleChunk(chunk)
	}

	returnChan := merge(results...)
	for i := range returnChan {
		result.charCount += i.charCount
		result.wordCount += i.wordCount
	}

	wg.Wait()
	executionTime := time.Since(startTime)
	writeResultFile(executionTime)
}

func writeResultFile(executionTime time.Duration) {
	mapString := ""
	for key, value := range result.frequency {
		mapString += fmt.Sprintf("[%s: %d]\n", key, value)
	}

	lines := []string{
		fmt.Sprintf("Number of lines: %d\n", result.lineCount),
		fmt.Sprintf("Number of words: %d\n", result.wordCount),
		fmt.Sprintf("Number of characters: %d\n", result.charCount),
		fmt.Sprintf("Average word length: %f\n", float64(result.totalWordLength)/float64(len(result.frequency))),
		fmt.Sprintln("Frequency:"),
		fmt.Sprintln(mapString),
	}

	current_date := time.Now().Format("02-01-2006")
	current_time := time.Now().Format("15-04-05")
	outputFileName := fmt.Sprintf("results/rs_%s_%s.txt", current_date, current_time)
	outputFile, err := os.Create(outputFileName)

	if err != nil {
		fmt.Println("error while writing")
		fmt.Println(err)
	}

	w := bufio.NewWriter(outputFile)

	for _, line := range lines {
		w.WriteString(line)
	}

	fmt.Printf("Number of lines: %d\n", result.lineCount)
	fmt.Printf("Number of words: %d\n", result.wordCount)
	fmt.Printf("Number of characters: %d\n", result.charCount)
	fmt.Printf("Average word length: %f\n", float64(result.totalWordLength)/float64(len(result.frequency)))
	fmt.Println("Frequency:")
	// fmt.Println(mapString)
	// fmt.Println(executionTime)
	fmt.Println("Execution Time:", executionTime)

	w.WriteString(fmt.Sprintln("Execution Time", executionTime))
	w.Flush()

	fmt.Println("Write file successfully")
	outputFile.Close()
}

func main() {
	startTime := time.Now()

	// filePath := "./documents/requirement_en.txt"
	// filePath := "./documents/file.txt"
	filePath := "./documents/file_1MB.txt"
	splitFileMethod(filePath, startTime)
}
