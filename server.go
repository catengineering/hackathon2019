package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
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

func (s Server) crash(w http.ResponseWriter, r *http.Request) {
	stringserver := instanceName()
	log.Printf(stringserver)
	serverName, err := os.Hostname()

	message := fmt.Sprintf("Instance %s Crashed.... ARGH The Zombies are coming.....", serverName)
	if err != nil {
		panic("cant determine hostname")
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, message)
	go func() {
		log.Printf("Crash\n")
		panic("ARGH... The Zombies are coming")
	}()

	return
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func instanceName() string {
	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2019-06-04", nil)
	check(err)
	req.Header.Set("Metadata", "true")
	resp, err := http.DefaultClient.Do(req)
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	return fmt.Sprintf(string(body))
}

func main() {

	var minDelay int64
	minDelay = 100
	var maxDelay int64
	maxDelay = 400

	if len(os.Args) > 2 {
		minDelayInput, err2 := strconv.ParseInt(os.Args[1], 10, 64)
		if (err2) != nil {
			log.Printf("Error parsing minDelayInput")
		}
		maxDelayInput, err3 := strconv.ParseInt(os.Args[2], 10, 64)
		if (err3) != nil {
			log.Printf("Error parsing maxDelayInput")
		}
		minDelay = int64(minDelayInput)
		maxDelay = int64(maxDelayInput)

	}
	delay, err := rand.Int(rand.Reader, big.NewInt(maxDelay-minDelay))
	if err != nil {
		panic(err)
	}
	delaySec := delay.Int64()
	delaySec += minDelay
	log.Printf("Delay %v (%v - %v)\n", delaySec, minDelay, maxDelay)
	crashesAt := time.Now().Add(time.Second * time.Duration(delaySec))

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
	http.HandleFunc("/crash", srv.crash)

	log.Println("Serving on port 8080")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal("ListenAndServe: ", err)
}
