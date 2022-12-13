package main

import (
	"fmt"
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
	go func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			superChanel <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	} ()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for {
			select {
			case data:=<-superChanel:
				tt, _ := time.Parse(time.RFC3339, data.cT)
				data.fT = time.Now().Format(time.RFC3339Nano)
				if tt.After(time.Now().Add(-20 * time.Second)) {
					data.taskRESULT = []byte("task has been successed")
					result[data.id] = data
				} else {
					data.taskRESULT = []byte("something went wrong")
					err = append(err, fmt.Errorf("Task id %d time %s, error %s", data.id, data.cT, data.taskRESULT))
				}
				time.Sleep(time.Millisecond * 150)
			case <-after:
				endChanel<-1
			}
		}
	}()
	<-endChanel

	println("Done tasks:")
	for dc := range result {
		fmt.Println(dc)
	}

	println("Errors:")
	for ec := range err {
		fmt.Println(ec)
	}

	fmt.Println("done working")
}