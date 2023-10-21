package golangweb

import (
	"fmt"
	"net/http"
	"testing"
)

// 2. Membuat server
func TestServer(t *testing.T) {
	server := http.Server{
		Addr: "localhost:8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// 3. handler
func TestHandler(t *testing.T) {

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		// w -> Writter atau yang akan dikirimkan ke client
		// r -> Request dari client

		// writter harus ditulis dengan []byte, tetapi dengan menggunakan Fprint()
		// maka konversi ke []byte sudah otomatis dilakukan.
		fmt.Fprint(w, "Hello World!")
	}

	server := http.Server{
		Addr: "localhost:8080",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

