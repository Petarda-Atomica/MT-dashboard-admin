package main

import "fmt"

func main() {
	key, err := retrieveKey([]byte("password"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(key))
}
