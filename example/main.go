package main

import (
	"github.com/MwlLj/go-session"
	_ "github.com/go-sql-driver/mysql"
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

func mysqlTest() {
	memory, err := session.New(&session.Memory_type_mysql)
	memory.Dial("root:123456@tcp(localhost:3306)/session")
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

func redisTest() {
	memory, err := session.New(&session.Memory_type_redis)
	memory.Dial("localhost:6379")
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
	var timeoutS int64 = 10
	for i := 0; i < 10; i++ {
		memory.Reset(id, &timeoutS)
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
	// memoryTest()
	// mysqlTest()
	redisTest()

	for {
		time.Sleep(10 * time.Second)
	}
}
