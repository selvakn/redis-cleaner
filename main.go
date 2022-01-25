package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	var connectionUrl = flag.String("c", "localhost:6379", "redis connection string")
	var redisPassword = flag.String("P", "", "redis connection password (optional)")
	var keyPrefix = flag.String("key-prefix", "", "key prefix to search (if not mention, then all keys)")
	var retainDaysMin = flag.Int("retain-days-min", 15, "keys to be retained with ttl less than this value")
	var retainDaysMax = flag.Int("retain-days-max", 450, "keys to be retained with ttl more than this value")

	flag.Parse()

	minHoursToRetain, err := time.ParseDuration(fmt.Sprintf("%dh", *retainDaysMin*24))
	if err != nil {
		panic("unable to parse retain days")
	}
	maxHoursToRetain, err := time.ParseDuration(fmt.Sprintf("%dh", *retainDaysMax*24))
	if err != nil {
		panic("unable to parse retain days")
	}

	fmt.Println("going to delete item between", minHoursToRetain, "and", maxHoursToRetain)

	rdb := redis.NewClient(&redis.Options{
		Addr:     *connectionUrl,
		Password: *redisPassword,
		DB:       0,
	})

	iter := rdb.Scan(ctx, 0, fmt.Sprintf("%s*", *keyPrefix), 0).Iterator()
	for iter.Next(ctx) {
		val := iter.Val()

		expiry, err := rdb.TTL(ctx, val).Result()
		if err != nil {
			fmt.Println("error on accessing ttl of key: ", val, err)
			continue
		}

		if expiry > minHoursToRetain && expiry < maxHoursToRetain {
			fmt.Println("key:", val, expiry, "deleting")

			if err := rdb.Del(ctx, val).Err(); err != nil {
				fmt.Println("error on deleting key: ", val, err)
				continue
			}

			continue
		}
		fmt.Println("key:", val, expiry, "retained")

	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
}
