package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andr-235/vk"
)

func main() {
	client := vk.New(
		vk.WithToken("YOUR_TOKEN"),
		vk.WithVersion("5.199"),
	)

	users, err := client.UsersGet(context.Background(), vk.UsersGetParams{
		UserIDs: []int{1},
		Fields:  []string{"photo_100"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", users)
}
