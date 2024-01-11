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
	remoteAddr := r.RemoteAddr
	log.Printf("request from %s", remoteAddr)

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error getting hostname: %s\n", err.Error())
		return
	}

	response := fmt.Sprintf("[%s] You've hit: %s\n", getTime(), hostname)
	fmt.Fprint(w, response)
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
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

	port := "8080"
	log.Printf("start listening on port http://127.0.0.1:%s\n", port)

	err := http.ListenAndServe(":"+PORT, nil)

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("server closed\n")
		} else {
			log.Fatal(err)
		}
	}
}
