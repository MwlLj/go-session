package main

import (
	"github.com/MwlLj/go-session"
	"log"
	"time"
)

func memoryTest() {
	memory, err := session.New(&session.Memory_type_memory)
	id, err := memory.Create(10)
	if err != nil {
		log.Println(err)
		return
	}
	if id != nil {
		log.Println(*id)
	}
	time.Sleep(1 * time.Second)
	// memory.Destory(id)
	for i := 0; i < 10; i++ {
		memory.Reset(id, nil)
		time.Sleep(1 * time.Second)
	}
	isValid, err := memory.IsValid(id)
	if err != nil {
		log.Println(err)
		return
	}
	if isValid {
		log.Println("is valid")
	}
}

func main() {
	memoryTest()

	for {
		time.Sleep(10 * time.Second)
	}
}
