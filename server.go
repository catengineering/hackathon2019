package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"encoding/json"
	"net/http"
)

type Server struct {
	ID        string
	CrashesAt time.Time
}

func (s Server) marcoPolo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var marco string

	if marcos, ok := r.Form["marco"]; ok {
		marco = marcos[0]
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Expected marco GET param"))
		return
	}

	w.WriteHeader(200)

	// (Doing some seriously hard work here)
	time.Sleep(time.Second / 2)
	log.Printf("marco: %s\n", marco)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"polo": marco,
		"id":   s.ID,
		"ttl":  fmt.Sprintf("%0.0f", s.CrashesAt.Sub(time.Now()).Seconds()),
	}); err != nil {
		return
	}
}

func (s Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func runtimeID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
	delay, err := rand.Int(rand.Reader, big.NewInt(900))
	if err != nil {
		panic(err)
	}

	crashesAt := time.Now().Add(time.Second * time.Duration(delay.Uint64()+100))

	srv := Server{
		ID:        runtimeID(),
		CrashesAt: crashesAt,
	}

	go func() {
		time.Sleep(srv.CrashesAt.Sub(time.Now()))
		log.Printf("Crashing\n")
		panic("oh no! what an unpredictable crash!")
	}()

	http.HandleFunc("/health", srv.health)
	http.HandleFunc("/", srv.marcoPolo)

	log.Println("Serving on port 8080")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal("ListenAndServe: ", err)
}
