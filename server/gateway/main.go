package main

import (
	"os"
	"log"
	"net/http"
	"github.com/capstone/uwcapstoneihmeportal.github.io/server/gateway/handlers"
	"github.com/nimajalali/go-force/force"
)

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

	consumerId := reqEnv("CONSUMER_ID")
	apiVersion := reqEnv("API_VERSION")
	consumerSecret := reqEnv("CONSUMER_SECRET")
	securityToken := reqEnv("SECURITY_TOKEN")
	forceUsername := reqEnv("FORCE_USERNAME")
	forcePassword := reqEnv("FORCE_PASSWORD")
	forceApiEnv := reqEnv("FORCE_API_ENV")

	//Create a forceApi Object
	forceApi, err := force.Create(apiVersion, consumerId, consumerSecret, forceUsername,
		forcePassword, securityToken, forceApiEnv)
	if err != nil {
		log.Fatal("Error creating force api with given credentials:", err)
	}

	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	ctx := handlers.NewContext(forceApi)

	tlsKeyPath := os.Getenv("TLSKEY")
	tlsCertPath := os.Getenv("TLSCERT")
	mux := http.NewServeMux()
	newMux := http.NewServeMux()
	newMux.Handle("/v1/authorize", handlers.NewCorsHandler(http.HandlerFunc(ctx.Authorize)))
	newMux.Handle("/v1/password_update", handlers.NewCorsHandler(http.HandlerFunc(ctx.PasswordUpdate)))
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
