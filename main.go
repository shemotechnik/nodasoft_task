package main

import (
	"fmt"
	"runtime"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChanel := make(chan Ttype, 10)
	endChanel := make(chan int)
	after := time.After(time.Second * 3)

	//create superChanel
	for i := 1; i <= runtime.GOMAXPROCS(0); i++ {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				superChanel <- Ttype{cT: ft, id: int(time.Now().UnixNano())} // передаем таск на выполнение
			}
		}()
	}

	doneChanel := make(chan Ttype, 10)
	errorChanel := make(chan error, 10)
	for i := 1; i <= runtime.GOMAXPROCS(0); i++ {
		go func() {
			for {
				select {
				case data := <-errorChanel:
					fmt.Printf("error task %s\n", data)
				case data := <-doneChanel:
					fmt.Printf("done task %d\n", data.id)
				case data := <-superChanel:
					tt, _ := time.Parse(time.RFC3339, data.cT)
					data.fT = time.Now().Format(time.RFC3339Nano)
					if tt.After(time.Now().Add(-20 * time.Second)) {
						data.taskRESULT = []byte("task has been successed")
						doneChanel <- data
					} else {
						data.taskRESULT = []byte("something went wrong")
						errorChanel <- fmt.Errorf("Task id %d time %s, error %s", data.id, data.cT, data.taskRESULT)
					}
					time.Sleep(time.Millisecond * 150)
				case <-after:
					endChanel <- 1
				}
			}
		}()
	}
	<-endChanel
}
