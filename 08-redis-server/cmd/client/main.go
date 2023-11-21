package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	maxReqs = 1000
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	light(client)
	stress(client)
}

func light(client *redis.Client) {
	st := client.Ping()
	if st.Err() != nil {
		panic(st.Err())
	}

	res, err := st.Result()
	if err != nil {
		panic(err)
	}

	println(res)

	str := client.Echo("Hello World")
	if str.Err() != nil {
		panic(str.Err())
	}

	res, err = str.Result()
	if err != nil {
		panic(err)
	}

	println(res)

	st = client.Set("key", "value", 0)
	if st.Err() != nil {
		panic(st.Err())
	}

	res, err = st.Result()
	if err != nil {
		panic(err)
	}

	println(res)

	str = client.Get("key")
	if str.Err() != nil {
		panic(st.Err())
	}

	res, err = str.Result()
	if err != nil {
		panic(err)
	}

	println(res)
}

func stress(client *redis.Client) {
	start := time.Now()
	for i := 1; i < maxReqs; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		st := client.Set(key, value, 0)
		if st.Err() != nil {
			panic(st.Err())
		}

		res, err := st.Result()
		if err != nil {
			panic(err)
		}
		if res != "OK" {
			panic("got wrong set response: " + res)
		}
	}

	fmt.Printf("%d requests took %s", maxReqs, time.Since(start).String())
}
