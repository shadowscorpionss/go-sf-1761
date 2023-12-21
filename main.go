package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var wg *sync.WaitGroup = &sync.WaitGroup{}

func integerGenerator(cap, pause int, ch <-chan int) <-chan int {
	immer := true
	c := make(chan int, cap)
	wg.Add(1)
	// Горутина-замыкание
	go func() {
		defer close(c)
		defer wg.Done()
		for immer {
			i := rand.Intn(10000)
			time.Sleep(time.Duration(pause) * time.Millisecond)
			c <- i
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch
		immer = false
	}()

	return c
}

func controlChannel() <-chan int {
	c := make(chan int)
	wg.Add(1)
	go func() {
		defer close(c)
		defer wg.Done()
		for {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if strings.ToLower(input) == "exit" {
				break
			}
		}

	}()
	return c
}

// Напишите код, в котором имеются два канала сообщений из целых чисел, так,
// чтобы приём сообщений из них никогда не приводил к блокировке
// и чтобы вероятность приёма сообщения из первого канала была выше в 2 раза, чем из второго.
// *Если хотите, можете написать код, который бы демонстрировал это соотношение.
// В качестве ответа приложите архивный файл с кодом программы из Задания 17.6.1.
func main() {
	fmt.Print("\033[2J") //Clear screen

	var ch1, ch2, ctrl <-chan int
	var s1, s2 int = 1, 1

	ctrl = controlChannel()
	ch1 = integerGenerator(200, 200, ctrl)
	ch2 = integerGenerator(100, 600, ctrl)

	wg.Add(1)
	read := func() {
		defer wg.Done()

	loop:
		for {

			fmt.Printf("\033[%d;%dH", 1, 0) // Set cursor position
			fmt.Printf("Принято сообщений на оба канала %d. Отношение принятых сообщений канал1/канал2 %d\n", s1+s2, s1/s2)

			select {
			case i1 := <-ch1:
				s1++
				fmt.Printf("\033[%d;%dH", 2, 0) // Set cursor position
				fmt.Printf("Канал 1: %d принято (%d %%). Последнее сообщение %d\n", s1, s1*100/(s1+s2), i1)
				fmt.Printf("\033[%d;%dH", 4, 0) // Set cursor position

			case i2 := <-ch2:
				s2++
				fmt.Printf("\033[%d;%dH", 3, 0) // Set cursor position
				fmt.Printf("Канал 2: %d принято (%d %%). Последнее сообщение %d\n", s2, s2*100/(s1+s2), i2)
				fmt.Printf("\033[%d;%dH", 4, 0) // Set cursor position

			case <-ctrl:
				break loop
			}

		}
	}

	go read()

	wg.Wait()
}
