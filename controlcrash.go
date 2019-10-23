package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"
	"os"
	//"bufio"
	//"io/ioutil"
	"encoding/json"
	"net/http"
)

type Server struct {
	ID        string
	CrashesAt time.Time
}
type CrashTimeRange struct {
	min,max int
}

func (s Server) marcoPolo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var marco string

	if marcos, ok := r.Form["peek"]; ok {
		marco = marcos[0]
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Expected peek GET param"))
		return
	}

	w.WriteHeader(200)

	// (Doing some seriously hard work here)
	time.Sleep(time.Second / 2)
	log.Printf("peek: %s\n", marco)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"a boo": marco,
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
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func (s Server) crash(w http.ResponseWriter, r *http.Request) {
	stringserver := instanceName()
	log.Printf(stringserver)
	serverName, err := os.Hostname()        
	
	message := fmt.Sprintf("Instance %s Crashed.... ARGH The Zombies are coming.....", serverName);
	if err != nil {
		panic("cant determine hostname")
	}
	
	w.WriteHeader(200)		
        	fmt.Fprintf(w,message);
		go func() {
	   		log.Printf("Crash\n")
	   		panic("ARGH... The Zombies are coming")	
		}()

	return
}

func instanceName() string {
	req, err := http.NewRequest("GET", "http://169.254.169.254/metadata/instance?api-version=2019-06-04", nil)
	check(err)
	req.Header.Set("Metadata", "true")
	resp, err := http.DefaultClient.Do(req)
	check(err)
	return fmt.Sprintf(resp.Body)
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
	delay, err := rand.Int(rand.Reader, big.NewInt(1700))
	if err != nil {
		panic(err)
	}

	crashesAt := time.Now().Add(time.Second * time.Duration(delay.Uint64()+100))

	srv := Server{
		ID:        runtimeID(),
		CrashesAt: crashesAt,
	}

	//go func() {
	//	time.Sleep(srv.CrashesAt.Sub(time.Now()))
	//	log.Printf("Crashing\n")
	//	panic("oh no! what an unpredictable crash!")
	//}()

	http.HandleFunc("/health", srv.health)
	http.HandleFunc("/", srv.marcoPolo)
	http.HandleFunc("/crash", srv.crash)
	log.Println("Serving on port 8080")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal("ListenAndServe: ", err)
}
