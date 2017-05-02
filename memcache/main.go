package main

import (
	"fmt"
  "github.com/go-redis/redis"
)
var client *redis.Client

func ExampleNewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "ec2-34-204-242-91.compute-1.amazonaws.com:29479",
		Password: "p258fce57afddd567a931a390a0694bad25fb6ca5f08d775a766d2c0d8f42a8ec", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {

	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exists
}

func main() {
	ExampleNewClient()
	//print(client)
	ExampleClient()
}
