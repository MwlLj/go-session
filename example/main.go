package main

import (
	"github.com/MwlLj/go-session"
	"log"
)

func memoryTest() {
	memory, err := session.New(&session.Memory_type_memory)
	id, err := memory.Create(10)
	if err != nil {
		log.Println(err)
	}
	if id != nil {
		log.Println(*id)
	}
}

func main() {
	memoryTest()
}
