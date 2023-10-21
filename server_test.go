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

// 3. handler (tidak mendukung multiple endpoint)
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

// 4. serveMux (support multiple endpoints)
func TestServeMux(t *testing.T) {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman >> /")
	})

	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman >> /hi")
	})

	// jika menambahkan / diakhir maka param apapun dibelakangnya akan
	// menggunakan handler ini. ex: /images/tidak-ada akan menggunakan handler ini
	mux.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman >> /images/")
	})

	// tetapi jika param berikutnya terdaftar pada handler, maka handlernya
	// akan digunakan. ex: /images/thumbnails/apapun akan menggunakan handler ini
	mux.HandleFunc("/images/thumbnails/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Halaman >> /images/thumbnails/")
	})

	// golang akan memprioritaskan url yang panjang terlebih daluhu.
	
	// maksudnya jika anda memasukkan url /images/thumbnails/ . kita tahu ada handler /images/
	// tetapi /images/thumbnails/ akan mencari handler /images/thumbnails/ terlebih dahulu.
	// jika /images/thumbnails/ ada maka akan digunakan jika tidak ada maka akan menggunakan
	// handler /images/ 

	server := http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	
}

