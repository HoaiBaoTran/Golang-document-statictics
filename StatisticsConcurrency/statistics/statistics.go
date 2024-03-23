package statistics

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal("Error while reading file", err)
	}
}

func CountFrequencyFromLine(data string) <-chan map[string]int {
	frequencyChan := make(chan map[string]int)
	go func() {
		frequency := make(map[string]int)
		totalWordLength := 0
		punctuationPattern := regexp.MustCompile(`[[:punct:]]`)
		cleanedString := punctuationPattern.ReplaceAllString(data, " ")
		words := strings.Fields(cleanedString)
		for _, word := range words {
			if _, isContain := frequency[word]; !isContain {
				totalWordLength += len(word)
			}
			frequency[word]++
		}
		frequencyChan <- frequency
	}()

	return frequencyChan
}

func WordAndCharCount(line string) (int, int) {
	wordSlice := strings.Fields(line)
	wordCount := len(wordSlice)

	lineWithoutSpace := strings.ReplaceAll(line, " ", "")
	charCount := len(lineWithoutSpace)

	return wordCount, charCount
}

func WriteResultToFile(lineCount, wordCount, charCount int, averageWordLength float64, frequency map[string]int, executionTime time.Duration) {
	mapString := ""
	for key, value := range frequency {
		mapString += fmt.Sprintf("[%s: %d]\n", key, value)
	}

	lines := []string{
		fmt.Sprintf("Number of lines: %d\n", lineCount),
		fmt.Sprintf("Number of words: %d\n", wordCount),
		fmt.Sprintf("Number of characters: %d\n", charCount),
		fmt.Sprintf("Average word length: %f\n", averageWordLength),
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

	fmt.Printf("Number of lines: %d\n", lineCount)
	fmt.Printf("Number of words: %d\n", wordCount)
	fmt.Printf("Number of characters: %d\n", charCount)
	fmt.Printf("Average word length: %f\n", averageWordLength)
	fmt.Println("Frequency:")
	// fmt.Println(result.frequency)
	// fmt.Println(mapString)
	fmt.Println("Execution Time:", executionTime)

	w.WriteString(fmt.Sprintln("Execution Time", executionTime))
	w.Flush()

	fmt.Println("Write file successfully")
	outputFile.Close()
}
