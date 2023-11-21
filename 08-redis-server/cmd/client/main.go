package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	ctx := context.Background()

	st := client.Ping(ctx)
	if st.Err() != nil {
		panic(st.Err())
	}

	res, err := st.Result()
	if err != nil {
		panic(err)
	}

	println(res)

	str := client.Echo(ctx, "Hello World")
	if str.Err() != nil {
		panic(str.Err())
	}

	res, err = str.Result()
	if err != nil {
		panic(err)
	}

	println(res)
}
