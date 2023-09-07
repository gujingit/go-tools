package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Movie struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// https://redis.io/docs/clients/go/
// https://www.golinuxcloud.com/go-crud-rest-api-redis-db/
func main() {
	opt, err := redis.ParseURL("redis://<user>:<pass>@localhost:6379/<db>")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)
	id := uuid.New().String()
	movie := Movie{
		Id:          id,
		Title:       "test",
		Description: "test description",
	}

	jsonBytes, err := json.Marshal(movie)
	if err != nil {
		panic(err)
	}
	client.HSet(context.TODO(), "movies", movie.Id, string(jsonBytes))
	if err != nil {
		panic(err)
	}

	val, err := client.HGet(context.TODO(), "movies", id).Result()
	if err != nil {
		panic(err)
	}
	m := &Movie{}
	err = json.Unmarshal([]byte(val), m)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", m)

}
