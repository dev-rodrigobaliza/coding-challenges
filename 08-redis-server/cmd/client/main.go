package main

import (
	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

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
