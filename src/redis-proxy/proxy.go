package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/karlseguin/ccache"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)
var max_connections = 10
var cache_size = 10
var cache_expiry = 60000
var redis_url = "redis:6379"
var port = "8080"
var redisPool *redis.Pool
var cache *ccache.Cache

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	log.Println("ABOUT")
	fmt.Fprintf(w, "Golang Redis proxy for Segment Interview by Ryan Nowacoski\n")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ROOT")
	http.Redirect(w, r, "https://github.com/rmn36/segment_redis_proxy", 301)
}

func createJSONOutput(status int, data string) string {
	outputMap := make(map[string]string)

	outputMap["status"] = strconv.Itoa(status)
	outputMap["data"] = data

	output, err := json.Marshal(outputMap)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	return string(output)
}

func performGet(key string) (int, string) {
	cacheResult := cache.Get(key)
	var getResult string
	if cacheResult != nil {
		log.Println("FOUND "+key+" IN CACHE")
		getResult = cacheResult.Value().(string)
	} else {
		c := redisPool.Get()
		defer c.Close()
	
		var err error
		getResult, err = redis.String(c.Do("GET", key))
		if err != nil {
			fmt.Printf("GET error: %v", err.Error())
			return -1, getResult
		}
	} 

	return 0, getResult
}

func performSet(key string, value string) (int) {
	c := redisPool.Get()
	defer c.Close()

	setResult, err := redis.String(c.Do("SET", key, value))
	if err != nil {
		fmt.Printf("SET error: %v\n", err.Error())
		return -1
	}

	if setResult == "OK" {
		cache.Set(key, value, time.Minute * time.Duration(cache_expiry))
		fmt.Printf("SET successful. '%v'\n", setResult)
		return 0
	} else {
		fmt.Printf("SET result was '%v' when SET '%v' '%v'.\n", setResult, key, value)
		return -1
	}
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()

	if len(r.Form["key"]) > 0 && len(r.Form["value"]) > 0 {
		key := r.Form["key"][0]
		value := r.Form["value"][0]
		log.Println("SET "+key+" "+value)

		status := performSet(key, value)

		fmt.Fprintf(w, createJSONOutput(status, ""))

	} else {
		fmt.Fprintf(w, createJSONOutput(-1, "Missing key and value parameters."))
	}
	
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()

	if len(r.Form["key"]) > 0 {
		key := r.Form["key"][0]
		log.Println("GET "+key)

		status, value := performGet(key)
		fmt.Fprintf(w, createJSONOutput(status, value))

	} else {
		fmt.Fprintf(w, createJSONOutput(-1, "Missing key parameter."))
	}
}

func readConfigs(){
	var conn_parse_err error
	max_connections, conn_parse_err := strconv.Atoi(os.Getenv("MAX_CONNECTIONS"))
	if conn_parse_err != nil {
		panic("MAX_CONNECTIONS MUST BE A NUMBER")
	}
	log.Println("MAX CONNECTIONS: "+strconv.Itoa(max_connections))

	var cache_parse_err error
	cache_size, cache_parse_err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if cache_parse_err != nil {
		panic("CACHE_SIZE MUST BE A NUMBER")
	}
	log.Println("CACHE SIZE: "+strconv.Itoa(cache_size))

	var cache_exp_parse_err error
	cache_expiry, cache_exp_parse_err = strconv.Atoi(os.Getenv("CACHE_EXPIRY_TIME"))
	if cache_exp_parse_err != nil {
		panic("CACHE_EXPIRY_TIME MUST BE A NUMBER")
	}
	log.Println("CACHE EXPIRY TIME (ms): "+strconv.Itoa(cache_expiry))

	redis_url = os.Getenv("REDIS_URL")
	log.Println("REDIS URL: "+redis_url)
	port = os.Getenv("PORT")
	log.Println("LISTENING PORT: "+port)
}

func createRedisPool() *redis.Pool {
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func () (redis.Conn, error) {
		  c, err := redis.Dial("tcp", redis_url)
		  if err != nil {
			return nil, err
		  }
		  return c, nil
		},
	  }

	pool.TestOnBorrow = func(c redis.Conn, t time.Time) error {
        if time.Since(t) < time.Minute {
            return nil
        }
        _, err := c.Do("PING")
        return err
    }

	pool.MaxActive = max_connections
	pool.IdleTimeout = time.Second * 10
	return pool
}

func main() {
	readConfigs()
	redisPool = createRedisPool()
	cache = ccache.New(ccache.Configure().MaxSize(int64(cache_size)).ItemsToPrune(1).GetsPerPromote(1))

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/about", aboutHandler)
	//http.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Server ended on port "+port)
}