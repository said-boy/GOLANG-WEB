package golangweb

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
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

		// cara mengambil parameter dari url.
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

// 7. mutliple parameter
func multipleParameterHandler(w http.ResponseWriter, r *http.Request) {
	nama := r.URL.Query().Get("nama")
	bahasa := r.URL.Query().Get("bahasa")

	// Fprintf -> print format 
	// %s -> string
	fmt.Fprintf(w, "Halo.. %s ", nama)
	fmt.Fprintf(w, "Kamu sedang belajar bahasa %s yaa..?", bahasa)
}

func TestMultipleParameter(t *testing.T) {
	w := httptest.NewRecorder()
	// untuk multiple gunakan & antara 1 parameter dan yang lainnya.
	r := httptest.NewRequest("GET", "http://localhost:8080/?nama=said&bahasa=golang", nil)
	multipleParameterHandler(w, r)

	fmt.Println(w.Body.String())
	
}

// 8. multiple values parameter
func HandlerMultipleValues(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	names := query["name"]

	fmt.Fprint(w, strings.Join(names, " "))
}

func TestMultipleValues(t *testing.T) {
	w := httptest.NewRecorder()
	
	// 1 parameter multi values dengan cara menggunakan key yang sama tetapi dengan value yang berbeda
	// dan tetap dipisahkan dengan tanda '&' untuk setiap query
	r := httptest.NewRequest("GET", "http://localhoat:8080/?name=muhammad&name=said&name=al&name=khudri", nil)

	HandlerMultipleValues(w,r)

	fmt.Println(w.Body.String())
}

// 9. Header 
func HandlerHeaderResponse(w http.ResponseWriter, r *http.Request){
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, r.Header.Get("content-type"))
}

func TestHeader(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

	HandlerHeaderResponse(w, r)
	fmt.Println(w.Header().Get("content-type"))
}

func HandlerHeaderRequest(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, r.Header.Get("x-powered-by"))

}

func TestHeaderRequest(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r.Header.Add("x-powered-by", "Muhammad Said Alkhudri")
	HandlerHeaderRequest(w, r)
	fmt.Println(w.Body.String())
}

func HandlerPostForm(w http.ResponseWriter, r *http.Request){
	// saat mengambil request post, golang harus melakukan parsing
	// sebelum diambil datanya.
	// err := r.ParseForm()
	// if err != nil {
	// 	panic(err)
	// }
	// f_name := r.PostForm.Get("first_name")
	// l_name := r.PostForm.Get("last_name")

	// tetapi anda dapat menggunakan postformvalue. ini sebenarnya
	// sebuah `shortcut` dari kode parsing diatas. sebenarnya
	// dibelakang layar yang dilakukan sama seperti diatas, yaitu 
	// diparsing terlebih dahulu.
	f_name := r.PostFormValue("first_name")
	l_name := r.PostFormValue("last_name")

	fmt.Fprintf(w, "Halo %s %s", f_name, l_name)
}

func TestPostForm(t *testing.T){
	formPost := strings.NewReader("first_name=Muhammad Said&last_name=Alkhudri")
	r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", formPost)
	w := httptest.NewRecorder()
	
	// saat mengirimkan post wajib menambahkan ini pada header requestnya
	// gunanya adalah untuk memberitahu server bahwa ini adalah data post
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	HandlerPostForm(w, r)
	fmt.Println(w.Body.String())
}

func HandlerResponseCode(w http.ResponseWriter, r *http.Request){
	token := r.Header.Get("Token")
	if token == "" {
		// WriteHeader digunakan untuk menambahkan status code
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w,"Silahkan login terlebih dahulu.")
	} else {
		// jika tidak ada WriteHeader maka default response nya 200 Ok
		fmt.Fprintln(w,"Selamat datang.")
	}
}

func TestResponseCodeInvalid(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	
	// Invalid Token
	r.Header.Add("Token", "")
	HandlerResponseCode(w, r)

	fmt.Println(w.Code)
	fmt.Println(w.Result().Status)
	fmt.Println(w.Body.String())
}

func TestResponseCodeValid(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	
	// Valid Token
	r.Header.Add("Token", "validToken")
	HandlerResponseCode(w, r)

	fmt.Println(w.Code)
	fmt.Println(w.Result().Status)
	fmt.Println(w.Body.String())
}

// Cookie
func HandlersetCookie(w http.ResponseWriter, r *http.Request){
	cookie := new(http.Cookie)
	cookie.Name = "X-Saidboy"
	cookie.Value = "Said"
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

func TestCookie(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	HandlersetCookie(w, r)
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		fmt.Println(cookie.Path)
	}
}

// cookie dari client
func cookieClient(w http.ResponseWriter, r *http.Request){
	cookie, err := r.Cookie("x-boy")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, cookie.Name)
	fmt.Fprint(w, cookie.Value)
}

func TestCookieClient(t *testing.T){
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	
	// mengirim cookie dari client
	cookie := &http.Cookie{
		Name: "x-boy",
		Value: "hallo boyyy",
	}
	r.AddCookie(cookie)
	
	// jangan lupa kirimkan ke fungsinya setelah ditambahkan.
	cookieClient(w, r)

	fmt.Println(w.Body.String())

}

// file server
// agar bisa mengakses file resource html kita
func TestFileServer(t *testing.T) {
	dir := http.Dir("./assets")
	file := http.FileServer(dir)

	mux := http.NewServeMux()
	// ingat -> /static/
	// jika tanpa http.StripPrefix maka harus (wajib) ada folder static
	// didalam folder assets, jika tidak maka akan 404.
	mux.Handle("/static/", file)

	// fungsi ini tidak dapat mengakses index.js

	server := http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// jika tidak ingin ada folder static didalam folder assets
// anda dapat menggunakan http.StripPrefix()
func TestFileServerStripPrefix(t *testing.T) {
	dir := http.Dir("./assets")
	file := http.FileServer(dir)

	mux := http.NewServeMux()
	// ingat -> /static/
	// maka static nya akan dihitung sebagai folder dan bukan suatu keharusan.
	mux.Handle("/static/", http.StripPrefix("/static",file))

	// fungsi ini dapat mengakses semuanya yang ada didalam
	// folder assets.

	server := http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// mengguanakan embed
//go:embed assets
var assets embed.FS

func TestFileServerEmbed(t *testing.T) {
	// secara otomatis akan masuk kedalam sub folder dari 
	// yang sudah ditentukan.
	dir , err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	file := http.FileServer(http.FS(dir))

	mux := http.NewServeMux()
	// ingat -> /static/
	// maka static nya akan dihitung sebagai folder dan bukan suatu keharusan.
	mux.Handle("/static/", http.StripPrefix("/static",file))

	// fungsi ini dapat mengakses semuanya yang ada didalam
	// folder assets.

	server := http.Server{
		Addr: "localhost:8080",
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}