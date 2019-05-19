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
	memory.Dial("localhost:6381")
	id, err := memory.Create(5000000000)
	if err != nil {
		log.Println(err)
		return
	}
	if id != nil {
		log.Println(*id)
	}
	time.Sleep(1 * time.Second)
	go memory.KeyTimeoutNtf()
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
	// extra param test
	extraInfos := make(map[string]string)
	extraInfos["userType"] = "admin"
	extraInfos["userUuid"] = "123456"
	id, err = memory.CreateWithMap(5000000000, &extraInfos)
	if err != nil {
		log.Printf("createWithMap error, err: %v\n", err)
		return
	}
	// get extra info
	isValid, extra, err := memory.IsValidWithMap(id)
	if err != nil {
		log.Printf("isValidWithMap error, err: %v\n", err)
		return
	}
	if isValid {
		for k, v := range *extra {
			log.Printf("key: %s, value: %s\n", k, v)
		}
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
