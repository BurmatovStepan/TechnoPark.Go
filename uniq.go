package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type UniqOptions struct {
	CountLines       bool
	DuplicatesOnly   bool
	UniqueOnly       bool
	IgnoreFields     int
	IgnoreCharacters int
	IgnoreCase       bool
}

func main() {
	options := initUniqOptions()
	args := flag.Args()

	var inputStream io.Reader
	var outputStream io.Writer

	switch len(args) {
	case 2:
		outputFile, err := os.Create(args[1])
		defer outputFile.Close()

		if err != nil {
			fmt.Println("An error occured while opening file " + args[1])
			return
		}

		outputStream = outputFile
		fallthrough

	case 1:
		inputFile, err := os.Open(args[0])
		defer inputFile.Close()

		if err != nil {
			fmt.Println("An error occured while opening file " + args[0])
			return
		}

		inputStream = inputFile
		if outputStream == nil {
			outputStream = os.Stdout
		}

	case 0:
		inputStream = os.Stdin
		outputStream = os.Stdout
	}

	var lines []string
	istream := bufio.NewScanner(inputStream)
	for istream.Scan() {
		lines = append(lines, istream.Text())
	}

	result, err := uniq(lines, options)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
		return
	}

	for _, line := range result {
		fmt.Fprintln(outputStream, line)
	}
}

func initUniqOptions() UniqOptions {
	var options UniqOptions = UniqOptions{}

	flag.BoolVar(&options.CountLines, "c", false, "Подсчитать количество встречаний строки во входных данных")
	flag.BoolVar(&options.DuplicatesOnly, "d", false, "Вывести только те строки, которые повторились во входных данных")
	flag.BoolVar(&options.UniqueOnly, "u", false, "Вывести только те строки, которые не повторились во входных данных")
	flag.IntVar(&options.IgnoreFields, "f", 0, "Не учитывать первые <num_fields> полей в строке")
	flag.IntVar(&options.IgnoreCharacters, "s", 0, "Не учитывать первые <num_chars> символов в строке")
	flag.BoolVar(&options.IgnoreCase, "i", false, "Не учитывать регистр")

	flag.Parse()

	return options
}

func uniq(lines []string, options UniqOptions) ([]string, error) {
	if len(lines) == 0 {
		return []string{}, nil
	}

	var optionsCount int = 0
	if options.CountLines {
		optionsCount++
	}
	if options.DuplicatesOnly {
		optionsCount++
	}
	if options.UniqueOnly {
		optionsCount++
	}

	if optionsCount > 1 {
		return []string{}, errors.New("Invalid options")
	}

	var result []string
	var currentLine string

	lineCount := 1
	currentUniqueLineIndex := 0
	currentUniqueLine := formatString(lines[0], options)
	for i := 1; i < len(lines); i++ {
		currentLine = formatString(lines[i], options)

		if currentUniqueLine != currentLine {
			result = appendResult(result, lines[currentUniqueLineIndex], lineCount, options)
			lineCount, currentUniqueLineIndex, currentUniqueLine = 0, i, currentLine
		}

		lineCount++
	}
	result = appendResult(result, lines[currentUniqueLineIndex], lineCount, options)

	return result, nil
}

func formatString(line string, options UniqOptions) string {
	const fieldSeparator = ' '
	result := line

	if options.IgnoreFields > 0 {
		fieldsCount := 0
		index := 0
		inField := false

		for i, character := range result {
			if character != fieldSeparator && !inField {
				fieldsCount++
				inField = true

				if fieldsCount > options.IgnoreFields {
					index = i
					break
				}

			} else if character == fieldSeparator {
				inField = false
			}

			index = i + 1
		}

		result = line[index:]
	}

	if options.IgnoreCharacters > 0 {
		runes := []rune(result)
		result = ""

		if options.IgnoreCharacters < len(runes) {
			result = string(runes[options.IgnoreCharacters:])
		}
	}

	if options.IgnoreCase {
		result = strings.ToLower(result)
	}

	return result
}

func appendResult(result []string, line string, lineCount int, options UniqOptions) []string {
	lineToAdd := line
	addLine := false

	switch {
	case options.DuplicatesOnly:
		if lineCount > 1 {
			lineToAdd = line
			addLine = true
		}

	case options.UniqueOnly:
		if lineCount == 1 {
			lineToAdd = line
			addLine = true
		}

	case options.CountLines:
		lineToAdd = fmt.Sprintf("%d %s", lineCount, line)
		addLine = true

	default:
		addLine = true
	}

	if addLine {
		result = append(result, lineToAdd)
	}

	return result
}
