package golangweb

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

	// lebih baik buat url se-unique mungkin. 

	server := http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	
}

// 5. Request
func TestRequest(t *testing.T) {

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.Method) // "GET"
		fmt.Fprintln(w, r.RequestURI) // "/"
		fmt.Fprintln(w, r.URL.Query().Get("nama")) // "/?nama=said" -> "said"
		fmt.Fprintln(w, r.Response) // <nil>
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

// 6. Http test
// Membuat implementasi seperti HandlerFunc
func myHandlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func TestHttpTest(t *testing.T) {

	// dengan menggunakan httptest kita tidak perlu lagi buka browser
	// untuk mengetahui hasil testnya.
	writter := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	myHandlerFunc(writter, request)

	// Versi singkat
	fmt.Println(writter.Body.String()) // "Hello World!"

	// versi lengkap
	response := writter.Result()
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body)) // "Hello World!"

}

