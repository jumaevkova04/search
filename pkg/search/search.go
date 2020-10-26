package search

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

// Result описывают один результат поиска.
type Result struct {
	// Фраза, которую искали
	Phrase string
	// Целиком вся строка, в котором нашли вхождение (без /n или /r/n в конце)
	Line string
	// Номер строки (начиная с 1), на которой нашли вхождение
	LineNum int64
	// Номер позиции (начиная с 1), на которой нашли вхождение
	ColNum int64
}

// All ищет все вхождения pharse в текстовых файлах files.
func All(ctx context.Context, phrаse string, files []string) <-chan []Result {
	ch := make(chan []Result)

	var j int
	for j = range files {
		j++
	}

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < j; i++ {
		wg.Add(1)
		go func(ctx context.Context, fileName string, ch chan<- []Result) {
			defer wg.Done()

			channel := FindAllPhraseInFile(phrаse, fileName)

			if len(channel) <= 0 {
				return
			}

			ch <- channel

		}(ctx, files[i], ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	cancel()
	return ch
}

// FindAllPhraseInFile ...
func FindAllPhraseInFile(phrase string, fileName string) []Result {
	result := []Result{}

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print(err)
		return result
	}

	text := string(file)

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		i++
		if strings.Contains(line, phrase) {
			r := Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}
			result = append(result, r)
		}
	}

	return result
}

// Any ищет любое одно вхождения pharse в текстовых файлах files.
func Any(ctx context.Context, phrаse string, files []string) <-chan Result {
	ch := make(chan Result)

	var j int
	for j = range files {
		j++
	}

	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < j; i++ {
		wg.Add(1)
		go func(ctx context.Context, fileName string, i int, ch chan<- Result) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				log.Printf("canceled %v", i)
				// return
			default:
				channel := FindAnyPhraseInFile(phrаse, fileName)

				if (Result{}) != channel {
					ch <- channel
				}
			}
		}(ctx, files[i], i, ch)
	}

	<-ch
	cancel()

	go func() {
		defer close(ch)
		wg.Wait()
		// cancel()
	}()

	return ch
}

// FindAnyPhraseInFile ...
func FindAnyPhraseInFile(phrase string, fileName string) (result Result) {
	// var result Result

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Print(err)
		return result
	}

	text := string(file)

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		i++
		if strings.Contains(line, phrase) {
			return Result{
				Phrase:  phrase,
				Line:    line,
				LineNum: int64(i),
				ColNum:  int64(strings.Index(line, phrase)) + 1,
			}
		}
	}

	return result
}
