package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var REDIS_KEY = [10]string{
	"neo",
	"trinity",
	"morpheus",
	"tank",
	"dozer",
	"mouse",
	"cypher",
	"apoc",
	"switch",
	"agent_smith",
}

const (
	REDIS_WRITE_TIME = 2
	REPORT_TIME      = 5
	REPORTS_COUNT    = 5
)

const REDIS_HOST = "localhost:6379"

const EMPTY_KEY_ERROR = "redis: nil"

type Action struct {
	Value int64 `json:"value"`
	Count int   `json:"count"`
}

func cleanRedis(ctx context.Context, client *redis.Client) {
	client.FlushAll(ctx)
}

func getRandomKey() string {
	return REDIS_KEY[rand.Intn(10)]
}

func getAction(ctx context.Context, client *redis.Client, actionKey string) (value Action, err error) {
	val, redisError := client.Get(ctx, actionKey).Result()

	if redisError != nil {
		switch redisError.Error() {
		case EMPTY_KEY_ERROR:
			value = Action{
				Value: 0,
				Count: 0,
			}

		default:
			err = redisError
		}
	}

	json.Unmarshal([]byte(val), &value)

	return
}

func redisStorer(ctx context.Context, client *redis.Client) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return

		default:
			var actionKey = getRandomKey()
			var timestamp = time.Now().UnixMilli()
			action, err := getAction(ctx, client, actionKey)

			if err != nil {
				log.Fatal("Error retrieving data from redis")
				return
			}

			action.Value = timestamp
			action.Count++

			byteAction, err := json.Marshal(action)

			if err != nil {
				return
			}

			setError := client.Set(ctx, actionKey, string(byteAction), 0).Err()

			if setError != nil {
				log.Fatal("Error storing data in redis")
				return
			}

			fmt.Printf("Storing value in redis Key: %s, Value: %v \n", actionKey, action)

			time.Sleep(REDIS_WRITE_TIME * time.Second)
		}
	}
}

func redisReporter(ctx context.Context, cancel context.CancelFunc, client *redis.Client, wg *sync.WaitGroup) {

	for reportsCount := 0; reportsCount < REPORTS_COUNT; reportsCount++ {
		printRetReport(ctx, client)

		time.Sleep(REPORT_TIME * time.Second)
	}

	wg.Done()
}

func printRetReport(ctx context.Context, client *redis.Client) (err error) {
	fmt.Println("")
	fmt.Println("Reporting values from redis")

	for _, key := range REDIS_KEY {
		action, getActionError := getAction(ctx, client, key)

		if err != nil {
			err = getActionError
			return
		}

		fmt.Printf("Key: %s, Value: %v \n", key, action)
	}

	return
}

func main() {
	var ctx = context.Background()
	var wg sync.WaitGroup
	var client = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithCancel(ctx)
	cleanRedis(ctx, client)

	wg.Add(1)

	go redisStorer(ctx, client)
	go redisReporter(ctx, cancel, client, &wg)

	wg.Wait()

	printRetReport(ctx, client)
	cancel()
	return
}
