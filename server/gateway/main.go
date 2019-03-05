package main

import (
	"os"
	"log"
	"net/http"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/go-redis/redis"
	"time"
	"net/http/httputil"
	"strings"
	"sync"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/handlers"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/models/users"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/indexes"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/sessions"
)

const headerUser = "X-User"
const maxConnRetries = 5

func GetUser(r *http.Request, store sessions.Store, signinKey string) (*users.User, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return nil, fmt.Errorf("Not Authenticated")
	}
	sessionState := &handlers.SessionState{}
	_, err := sessions.GetState(r, signinKey, store, sessionState);
	if err != nil {
		return nil, fmt.Errorf("error getting state %v", err)
	}
	return sessionState.AuthUser, nil
}

func NewServiceProxy(addrs string, store sessions.Store, signinKey string) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	addrIndex := 0
	mx := sync.Mutex{}
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[addrIndex]
			addrIndex = (addrIndex + 1) % len(splitAddrs)
			mx.Unlock()

			r.Header.Del(headerUser)
			user, err := GetUser(r, store, signinKey)
			log.Println(user)
			if err != nil {
				return
			}
			userJSON, _ := json.Marshal(user)
			r.Header.Set(headerUser, string(userJSON))
		},
	}
}
//main is the main entry point for the server
func main() {
	/* TODO: add code to do the following
	- Create a new mux for the web server.
	- Tell the mux to call your handlers.SummaryHandler function
	  when the "/v1/summary" URL path is requested.
	- Start a web server listening on the address you read from
	  the environment variable, using the mux you created as
	  the root handler. Use log.Fatal() to report any errors
	  that occur when trying to start the web server.
	*/
	//get ADDR env variable

	mysqlAddr := reqEnv("MYSQL_ADDR")
	mysqlDB := reqEnv("MYSQL_DATABASE")
	mysqlPwd := reqEnv("MYSQL_ROOT_PASSWORD")
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}
	config := mysql.Config{
		Addr: mysqlAddr,
		User: "root",
		Passwd: mysqlPwd,
		DBName: mysqlDB,
		Net: "tcp",
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Fatalf("error opening database: %v\n", err)
	}
	defer db.Close()
	userStore := users.NewMySQLStore(db)
	lookUpTrie := indexes.NewTrie()
	err = userStore.LoadTrie(lookUpTrie)
	if err != nil {
		log.Fatalf("error loading users into trie: %v", err)
	}
	redisaddr := reqEnv("REDISADDR")
	signinKey := reqEnv("SESSIONKEY")
	summaryAddrs := reqEnv("SUMMARYADDRS")
	messageAddrs := reqEnv("MESSAGESADDR")

	if len(redisaddr) == 0 {
		redisaddr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisaddr,
	})
	store := sessions.NewRedisStore(client, time.Hour)
	notifier := handlers.NewNotifier()
	ctx := handlers.NewContext(signinKey, store, userStore, lookUpTrie, notifier)

	mqAddr := reqEnv("MQADDR")
	if len(mqAddr) == 0 {
		mqAddr = "127.0.0.1:5672"
	}
	conn, err := connectToMQ(mqAddr)
	if err != nil {
		log.Fatalf("error dialing MQ: %v", err)
	}

	handlers.RabbitMQSubscriber(notifier, conn)

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	mux := http.NewServeMux()
	newMux := http.NewServeMux()
	newMux.Handle("/v1/users", handlers.NewCorsHandler(http.HandlerFunc(ctx.UsersHandler)))
	newMux.Handle("/v1/users/", handlers.NewCorsHandler(http.HandlerFunc(ctx.SpecificUserHandler)))
	newMux.Handle("/v1/sessions", handlers.NewCorsHandler(http.HandlerFunc(ctx.SessionsHandler)))
	newMux.Handle("/v1/sessions/", handlers.NewCorsHandler(http.HandlerFunc(ctx.SpecificSessionHandler)))

	newMux.Handle("/v1/users/me/starred/messages/", handlers.NewCorsHandler(NewServiceProxy(messageAddrs, store, signinKey)))
	newMux.Handle("/v1/users/me/starred/messages", handlers.NewCorsHandler(NewServiceProxy(messageAddrs, store, signinKey)))
	newMux.Handle("/v1/channels/", handlers.NewCorsHandler(NewServiceProxy(messageAddrs, store, signinKey)))
	newMux.Handle("/v1/channels", handlers.NewCorsHandler(NewServiceProxy(messageAddrs, store, signinKey)))
	newMux.Handle("/v1/messages/", handlers.NewCorsHandler(NewServiceProxy(messageAddrs, store, signinKey)))
	newMux.Handle("/v1/summary", handlers.NewCorsHandler(NewServiceProxy(summaryAddrs, store, signinKey)))
	newMux.Handle("/v1/ws", handlers.NewCorsHandler(http.HandlerFunc(ctx.WebsocketUpgradeHandler)))


	mux.Handle("/v1/", newMux)
	log.Printf("server is listening at https://%s", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, mux))
}

func reqEnv(name string) string {
	val := os.Getenv(name)
	if len(val) == 0{
		log.Fatalf("please set %s variable", name)
	}
	return val
}

func connectToMQ(addr string) (*amqp.Connection, error) {
	mqURL := "amqp://" +  addr
	var conn *amqp.Connection
	var err error
	for i := 1; i <= maxConnRetries; i++ {
		conn, err = amqp.Dial(mqURL)
		if err == nil {
			log.Printf("successfully connected to MQ")
			return conn, nil
		}
		log.Printf("error connecting to MQ at %s: %v", mqURL, err)
		log.Printf("will retry in %d seconds", i*2)
		time.Sleep(time.Second * time.Duration(i * 2))
	}
	return nil, err
}
