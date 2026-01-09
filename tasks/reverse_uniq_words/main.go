package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type pair struct {
	word         string
	reversedWord string
}

// начало решения

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string, 5)
	go func() {
		defer close(out)
		for {
			select {
			case out <- randomWord(5):
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string, 5)
	go func() {
		defer close(out)
		for {
			select {
			case word, ok := <-in:
				if !ok {
					return
				}
				seen := make(map[rune]bool)
				flag := true
				for _, letter := range word {
					if seen[letter] {
						flag = false
						break
					}
					seen[letter] = true
				}
				if flag {
					select {
					case out <- word:
					case <-cancel:
						return
					}
				}
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan pair {
	out := make(chan pair, 5)
	go func() {
		defer close(out)
		for {
			select {
			case word, ok := <-in:
				if !ok {
					return
				}
				reversed := []rune(word)
				for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
					reversed[i], reversed[j] = reversed[j], reversed[i]
				}
				select {
				case out <- pair{string(reversed), word}:
				case <-cancel:
					return
				}
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan pair) <-chan pair {
	out := make(chan pair, 5)
	wg := sync.WaitGroup{}
	wg.Go(func() {

		for words := range c1 {
			select {
			case out <- words:
			case <-cancel:
				return
			}
		}
	})
	wg.Go(func() {
		for words := range c2 {
			select {
			case out <- words:
			case <-cancel:
				return
			}
		}
	})

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan pair, n int) {
	go func() {
		for i := 0; i < n; i++ {
			select {
			case word, ok := <-in:
				if !ok {
					return
				}
				fmt.Println(word.word, word.reversedWord)
			case <-cancel:
				return
			}
		}
	}()
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
