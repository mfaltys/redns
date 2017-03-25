package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixvoid/glogger"
	"gopkg.in/redis.v5"
)

var (
	loglevel  = "debug"
	redisHost = "localhost:6379"
	redisPass = ""
)

func main() {
	// TODO
	//   - flags for overriding: redishost, redispass, loglevel

	// initialize logger
	initLogger(loglevel)

	// initialize redis connection
	redisClient, err := initRedisConnection()
	if err != nil {
		glogger.Debug.Println("redis conneciton cannot be made, trying again in 1 second")
		redisClient, err = initRedisConnection()
		if err != nil {
			glogger.Error.Println("redis connection cannot be made.")
			os.Exit(1)
		}
	} else {
		glogger.Debug.Println("connection to redis succeeded.")
		glogger.Info.Println("link to redis on", redisHost)
	}

	// read in args
	args := os.Args[1:]

	if len(args) == 0 {
		// nothing was entered, end
		glogger.Error.Println("no arguments passed.")
		os.Exit(0)
	}

	switch args[0] {
	case "list":
		// list clients
		glogger.Debug.Println("listing clients")
		println()
		listClients(redisClient)
	case "get":
		// check if a client name is actually given
		if len(args) != 2 {
			// nothing was entered, end
			glogger.Error.Println("no client given.")
			// TODO : print syntax of ./doic_cli get <client_ip>
			os.Exit(0)
		}
		// get client history for args[1]
		glogger.Debug.Println("getting client entries")
		println()
		getClientHistory(redisClient, args[1])
	}
}

func initLogger(logLevel string) {
	// init logger
	if logLevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if logLevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if logLevel == "info" {
		glogger.LogInit(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		glogger.LogInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	}
}

func initRedisConnection() (*redis.Client, error) {
	// initialize redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPass,
		DB:       0,
	})
	_, redisErr := client.Ping().Result()
	return client, redisErr
}

func listClients(redisClient *redis.Client) {
	// get client:list from redis db
	clientList, err := redisClient.SMembers("client:list").Result()
	if err != nil {
		glogger.Error.Printf("error while getting 'client:list': %s", err)
	}

	println("CLIENTS:")
	fmt.Printf("%s\n", clientList)
	//glogger.Debug.Printf("%s", clientList)
}

func getClientHistory(redisClient *redis.Client, client string) {
	// get client:list from redis db
	clientHistory, err := redisClient.SMembers(fmt.Sprintf("client:%s", client)).Result()
	if err != nil {
		glogger.Error.Printf("error while getting '%s' history: %s", client, err)
	}

	fmt.Printf("HISTORY for '%s':\n", client)
	fmt.Printf("%s\n", clientHistory)
	//glogger.Debug.Printf("%s", clientHistory)
}
