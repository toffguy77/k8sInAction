package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
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
	slog.Info("requst got from host:", slog.String("remoteAddr", remoteAddr))

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("get hostname failed", slog.String("remoteAddr", remoteAddr), slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("[%s] You've hit: %s\n", getTime(), hostname)
	slog.Info("response", slog.String("response", response))
	fmt.Fprint(w, response)
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func checkLiveness(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(10)%2 == 0 {
		slog.Error("liveness check failed", slog.String("remoteAddr", r.RemoteAddr))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Not ok!")
		return
	}
	slog.Info("liveness check ok", slog.String("remoteAddr", r.RemoteAddr))
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Ok")
}

func main() {
	http.HandleFunc("/", getHost)
	http.HandleFunc("/liveness", checkLiveness)

	port := "8080"
	slog.Info("start listening on port http://127.0.0.1:%s", slog.String("port", port))

	err := http.ListenAndServe(":"+PORT, nil)

	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("server closed", slog.String("port", port))
		} else {
			slog.Error("server error", slog.String("port", port), slog.String("error", err.Error()))
		}
	}
}
