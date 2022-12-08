package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type InputFileReader interface {
	io.ReadCloser
	ReadLine() (line []byte, isPrefix bool, err error)
	ReadAll(callback func(string) error) error
}

type inputFileReader struct {
	file      *os.File
	bufreader *bufio.Reader
}

func OpenInputFile(name string) (InputFileReader, error) {
	filename := fmt.Sprintf("./inputs/%s", name)
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	ifr := inputFileReader{
		file:      file,
		bufreader: nil,
	}
	return &ifr, nil
}

func OpenAndReadAll(name string, callback func(string) error) error {
	ifr, err := OpenInputFile(name)
	if err != nil {
		return err
	}
	err = ifr.ReadAll(callback)
	if err != nil {
		return err
	}
	err = ifr.Close()
	if err != nil {
		return err
	}
	return nil
}

func OpenAndReadAllRegex(name string, regex string, allMustMatch bool) ([][]string, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	var results [][]string
	err = OpenAndReadAll(name, func(s string) error {
		m := re.FindStringSubmatch(s)
		if m == nil {
			if allMustMatch {
				return fmt.Errorf("non-matching line: %s", s)
			}
		} else {
			results = append(results, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func StringsToInts(s []string, positions ...int) ([]int, error) {
	var results []int
	for _, p := range positions {
		v, err := strconv.Atoi(s[p])
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

func (ifr *inputFileReader) Read(p []byte) (n int, err error) {
	return ifr.bufreader.Read(p)
}

func (ifr *inputFileReader) Close() error {
	return ifr.file.Close()
}

func (ifr *inputFileReader) ReadLine() (line []byte, isPrefix bool, err error) {
	if ifr.bufreader == nil {
		ifr.bufreader = bufio.NewReader(ifr.file)
	}
	return ifr.bufreader.ReadLine()
}

func (ifr *inputFileReader) ReadAll(callback func(string) error) error {
	scanner := bufio.NewScanner(ifr.file)
	var err error
	for scanner.Scan() {
		err = callback(scanner.Text())
		if err != nil {
			break
		}
	}
	if err == nil {
		err = scanner.Err()
	}
	return err
}
