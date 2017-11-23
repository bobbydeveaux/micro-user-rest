package main

import (
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"os"
	"time"
)

type person struct {
	Id          int64
	Name        string
	Valid       bool
	Jwt         string
	AccessToken string
}

func main() {

	router := mux.NewRouter()
	http.Handle("/", httpInterceptor(router))

	//router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/login", userLogin).Methods("GET")

	http.ListenAndServe(":8181", nil)

}

func httpInterceptor(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if origin := req.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}
		//startTime := time.Now()

		router.ServeHTTP(w, req)

		//finishTime := time.Now()
		//elapsedTime := finishTime.Sub(startTime)

		switch req.Method {
		case "GET":
			// We may not always want to StatusOK, but for the sake of
			// this example we will
		//	LogAccess(w, req, elapsedTime)
		case "POST":
			// here we might use http.StatusCreated
		}

	})
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	//NATS_HOST = nats://localhost:4222
	nc, _ := nats.Connect(os.Getenv("NATS_HOST"))
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Stranger"
	}

	me := &person{
		Name: name,
	}

	// Go type Request
	var p person
	err := ec.Request("user.login", me, &p, 1000*time.Millisecond)
	if err != nil {
		if nc.LastError() != nil {
			log.Println("Error in Request: %v\n", nc.LastError())
		}
		log.Println("Error in Request: %v\n", err)
	} else {
		log.Printf("Published [%s] : '%s'\n", "user.login", p.Name)
		log.Printf("Received User [%v] : '%s'\n", p.Id, p.Name)
	}

	b, err := json.Marshal(p)
	if err != nil {
		log.Println("error:", err)
	}

	if b == nil {
		b = []byte("{\"error\":\"timeout\"}")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}
