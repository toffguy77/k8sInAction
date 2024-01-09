package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	PORT = "8080"
)

func getHost(w http.ResponseWriter, r *http.Request) {
	log.Printf("request from %s", r.RemoteAddr)
	hn, err := os.Hostname()
	if err != nil {
		io.WriteString(w, "Error getting hostname: "+err.Error()+"\n")
	}
	io.WriteString(w, "["+getTime()+"] "+"You've hit: "+hn+"\n")
}

func getTime() string {
	t := time.Now()
	return fmt.Sprint(t.Format("2006-01-02 15:04:05"))
}

func checkLiveness(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(10)%2 == 0 {
		log.Printf("liveness check failed...")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Not ok!")
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Ok")
}

func main() {
	http.HandleFunc("/", getHost)
	http.HandleFunc("/liveness", checkLiveness)

	fmt.Printf("start lintening on port http://127.0.0.1:%s\n", PORT)
	err := http.ListenAndServe(":"+PORT, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		panic(err)
	}
}
