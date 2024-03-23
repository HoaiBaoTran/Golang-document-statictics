package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type Result struct {
	lineCount         int
	wordCount         int
	charCount         int
	averageWordLength float64
	frequency         map[string]int
}

func writeResultFile(result Result, executionTime time.Duration) {
	mapString := ""
	for key, value := range result.frequency {
		mapString += fmt.Sprintf("[%s: %d]\n", key, value)
	}

	lines := []string{
		fmt.Sprintf("Number of lines: %d\n", result.lineCount),
		fmt.Sprintf("Number of words: %d\n", result.wordCount),
		fmt.Sprintf("Number of characters: %d\n", result.charCount),
		fmt.Sprintf("Average word length: %f\n", result.averageWordLength),
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
	fmt.Printf("Average word length: %f\n", result.averageWordLength)
	// fmt.Println("Frequency:")
	// fmt.Println(result.frequency)
	fmt.Println("Execution Time:", executionTime)

	w.WriteString(fmt.Sprintln("Execution Time", executionTime))
	w.Flush()

	fmt.Println("Write file successfully")
	outputFile.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Error while reading file")
	}
}

func countWordAndCharFromLine(line string) (int, int) {
	words := strings.Fields(line)
	removeSpaceLine := strings.ReplaceAll(line, " ", "")
	return len(words), len(removeSpaceLine)
}

func countFrequencyFromLine(line string, frequency map[string]int) map[string]int {
	regexPattern := `(\w+)`
	re := regexp.MustCompile(regexPattern)
	words := strings.Fields(line)
	for _, word := range words {
		formattedWords := re.FindStringSubmatch(word)
		if len(formattedWords) < 1 {
			continue
		}
		frequency[formattedWords[0]]++
	}
	return frequency
}

func calculateAverageWordLength(frequency map[string]int) float64 {
	totalCharLength := 0
	for key := range frequency {
		totalCharLength += len(key)
	}
	return float64(totalCharLength) / float64(len(frequency))
}

func statisticsDocs(filePath string, startTime time.Time) {

	result := Result{
		frequency: make(map[string]int),
	}

	requirementFile, err := os.Open(filePath)
	checkError(err)

	scanner := bufio.NewScanner(requirementFile)

	for scanner.Scan() {
		line := scanner.Text()
		wordCount, charCount := countWordAndCharFromLine(line)
		result.wordCount += wordCount
		result.charCount += charCount
		result.frequency = countFrequencyFromLine(line, result.frequency)
		result.lineCount++
	}

	requirementFile.Close()

	result.averageWordLength = calculateAverageWordLength(result.frequency)

	executionTime := time.Since(startTime)
	writeResultFile(result, executionTime)
}

func main() {
	startTime := time.Now()
	// filePath := "../documents/requirement_en.txt"
	// filePath := "../documents/file.txt"
	filePath := "../documents/file_1MB.txt"
	statisticsDocs(filePath, startTime)
}
