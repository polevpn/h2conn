package main

import (
	"log"
	"net/http"

	"github.com/polevpn/h2conn"
)

func main() {
	srv := &http.Server{Addr: ":8000", Handler: handler{}}
	log.Printf("Serving on http://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h2conn.Accept(w, r)
	if err != nil {
		log.Printf("Failed creating connection from %s: %s", r.RemoteAddr, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	// Conn has a RemoteAddr property which helps us identify the client
	log.Println(conn.RemoteAddr().String(), "Joined")
	defer log.Printf("Left")

	// Loop forever until the client hangs the connection, in which there will be an error
	// in the decode or encode stages.
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Failed decoding request: %v", err)
			return
		}

		log.Printf("Got: %v", string(buf[:n]))
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Failed encoding response: %v", err)
			return
		}
	}
}
