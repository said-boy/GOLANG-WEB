package golangweb

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"text/template"
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
		Addr:    "localhost:8080",
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
		Addr:    "localhost:8080",
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
		fmt.Fprintln(w, r.Method)     // "GET"
		fmt.Fprintln(w, r.RequestURI) // "/"

		// cara mengambil parameter dari url.
		fmt.Fprintln(w, r.URL.Query().Get("nama")) // "/?nama=said" -> "said"
		fmt.Fprintln(w, r.Response)                // <nil>
	}

	server := http.Server{
		Addr:    "localhost:8080",
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

	HandlerMultipleValues(w, r)

	fmt.Println(w.Body.String())
}

// 9. Header
func HandlerHeaderResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	fmt.Fprint(w, r.Header.Get("content-type"))
}

func TestHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

	HandlerHeaderResponse(w, r)
	fmt.Println(w.Header().Get("content-type"))
}

func HandlerHeaderRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Header.Get("x-powered-by"))

}

func TestHeaderRequest(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	r.Header.Add("x-powered-by", "Muhammad Said Alkhudri")
	HandlerHeaderRequest(w, r)
	fmt.Println(w.Body.String())
}

func HandlerPostForm(w http.ResponseWriter, r *http.Request) {
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

func TestPostForm(t *testing.T) {
	formPost := strings.NewReader("first_name=Muhammad Said&last_name=Alkhudri")
	r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/", formPost)
	w := httptest.NewRecorder()

	// saat mengirimkan post wajib menambahkan ini pada header requestnya
	// gunanya adalah untuk memberitahu server bahwa ini adalah data post
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	HandlerPostForm(w, r)
	fmt.Println(w.Body.String())
}

func HandlerResponseCode(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token == "" {
		// WriteHeader digunakan untuk menambahkan status code
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Silahkan login terlebih dahulu.")
	} else {
		// jika tidak ada WriteHeader maka default response nya 200 Ok
		fmt.Fprintln(w, "Selamat datang.")
	}
}

func TestResponseCodeInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:8080/", nil)

	// Invalid Token
	r.Header.Add("Token", "")
	HandlerResponseCode(w, r)

	fmt.Println(w.Code)
	fmt.Println(w.Result().Status)
	fmt.Println(w.Body.String())
}

func TestResponseCodeValid(t *testing.T) {
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
func HandlersetCookie(w http.ResponseWriter, r *http.Request) {
	cookie := new(http.Cookie)
	cookie.Name = "X-Saidboy"
	cookie.Value = "Said"
	cookie.Path = "/"
	http.SetCookie(w, cookie)
}

func TestCookie(t *testing.T) {
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
func cookieClient(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("x-boy")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, cookie.Name)
	fmt.Fprint(w, cookie.Value)
}

func TestCookieClient(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)

	// mengirim cookie dari client
	cookie := &http.Cookie{
		Name:  "x-boy",
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
		Addr:    "localhost:8080",
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
	mux.Handle("/static/", http.StripPrefix("/static", file))

	// fungsi ini dapat mengakses semuanya yang ada didalam
	// folder assets.

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// mengguanakan embed
//
//go:embed assets
var assets embed.FS

func TestFileServerEmbed(t *testing.T) {
	// secara otomatis akan masuk kedalam sub folder dari
	// yang sudah ditentukan.
	dir, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	file := http.FileServer(http.FS(dir))

	mux := http.NewServeMux()
	// ingat -> /static/
	// maka static nya akan dihitung sebagai folder dan bukan suatu keharusan.
	mux.Handle("/static/", http.StripPrefix("/static", file))

	// fungsi ini dapat mengakses semuanya yang ada didalam
	// folder assets.

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// mengirim file dari server
// http.ServeFile
func fileServer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("nama") != "" {
		http.ServeFile(w, r, "pages/status/ok.html")
	} else {
		http.ServeFile(w, r, "pages/status/404.html")
	}
}

func TestFileFromServer(t *testing.T) {
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(fileServer),
	}
	server.ListenAndServe()
}

// menggunakan embed
//
//go:embed pages/status/404.html
var notFound embed.FS

// jika menggunakan string, pastikan data yang akan
// diembed sudah pasti berupa string
//
//go:embed pages/status/ok.html
var ok string

func fileServerEmbed(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("nama") != "" {
		fmt.Fprint(w, ok)
	} else {
		file, _ := notFound.ReadFile("pages/status/404.html")
		fmt.Fprint(w, string(file))
	}
}

func TestFileFromServerEmbed(t *testing.T) {
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(fileServerEmbed),
	}
	server.ListenAndServe()
}

// ParseFile()
// memanggil template html (tetapi cuma 1)
func Template(w http.ResponseWriter, r *http.Request) {
	// ParseFile -> untuk mengambil template dari file .html
	// harus sepesifik file.
	t := template.Must(template.ParseFiles("pages/index.gohtml"))
	err := t.ExecuteTemplate(w, "index.gohtml", "Halo")
	if err != nil {
		panic(err)
	}
}

func TestTemplate(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	Template(w, r)
	fmt.Println(w.Body.String())
}

// ParseGlob()
// memanggil file html semuanya dengan tanda *
func TemplateParseGlob(w http.ResponseWriter, r *http.Request) {
	// ParseGlob -> untuk dapat memanggil semua template.
	t := template.Must(template.ParseGlob("pages/*.gohtml"))
	err := t.ExecuteTemplate(w, "index.gohtml", "Halo")
	if err != nil {
		panic(err)
	}
}

func TestTemplateParseGlob(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	TemplateParseGlob(w, r)
	fmt.Println(w.Body.String())
}

type Address struct {
	Street string
}

type Index struct {
	Title   string
	Name    string
	Address Address
}

// mengirim data ke template
func TemplateData(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("templates/*.gohtml"))
	err := t.ExecuteTemplate(w, "index.gohtml", Index{
		Title: "Hal Index Ini boy",
		Name:  "Muhammad Said Alkhudri",

		// data bersarang
		Address: Address{
			Street: "Jl. H. Usman",
		},
	})
	if err != nil {
		panic(err)
	}
}

func TestTemplateData(t *testing.T) {
	// w := httptest.NewRecorder()
	// r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	// TemplateData(w, r)
	// fmt.Println(w.Body.String())

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: http.HandlerFunc(TemplateData),
	}
	server.ListenAndServe()
}

// mengirim data ke template
func TemplateAction(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseGlob("templates/*.gohtml"))
	err := t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title": "Hal Index Ini boy",
		// Name: "Muhammad Said Alkhudri",

		// data bersarang
		"Address": Address{
			Street: "Jl. H. Usman",
		},
		"Hobies": []string{
			"Makan",
			"Ngoding",
		},
		"Score": 90,
	})
	if err != nil {
		panic(err)
	}
}

func TestTemplateAction(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	TemplateAction(w, r)
	fmt.Println(w.Body.String())
}

// template function
func TemplateFunction(w http.ResponseWriter, r *http.Request) {
	t := template.New("function")
	t.Funcs(template.FuncMap{
		"upper": func(name string) string {
			return strings.ToUpper(name)
		},
	})
	k := template.Must(t.ParseFiles("templates/function.html"))
	k.ExecuteTemplate(w, "function", map[string]interface{}{
		"Name": "said",
	})
}

func TestTemplateFunction(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	TemplateFunction(w, r)
	fmt.Println(w.Body.String())
}

// caching template
// agar tidak membuat terus saat membutuhkan handler

var templateCaching = template.Must(template.ParseFiles("templates/templateCaching.html"))

func TemplateCaching(w http.ResponseWriter, r *http.Request) {
	templateCaching.ExecuteTemplate(w, "templateCaching.html", map[string]interface{}{
		"Name": "Said",
	})
}

func TestTemplateCaching(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	TemplateCaching(w, r)
	fmt.Println(w.Body.String())
}

// pipeline
// meneruskan return function ke function selanjutnya
var templatePipeline = template.New("templatePipeline.html")

func TemplatePipeline(w http.ResponseWriter, r *http.Request) {
	templatePipeline.Funcs(template.FuncMap{
		"upper": func(name string) string {
			return strings.ToUpper(name)
		},
		"count": func(name string) int {
			return strings.Count(name, "U")
		},
	})
	templatePipeline.ParseFiles("templates/templatePipeline.html")
	templatePipeline.ExecuteTemplate(w, "templatePipeline.html", map[string]interface{}{
		"Name": "Muhammad Said Alkhudri",
	})
}

func TestTemplatePipeline(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
	TemplatePipeline(w, r)
	fmt.Println(w.Body.String())
}

// redirect
func redirectFrom(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/redirect-to", http.StatusPermanentRedirect)
}

func redirectTo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Halo world")
}

func TestRedirect(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/redirect-from", redirectFrom)
	mux.HandleFunc("/redirect-to", redirectTo)

	// w := httptest.NewRecorder()
	// r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/redirect-from", nil)
	// redirectTo(w, r)

	// fmt.Println(w.Body.String())

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()

	// penggunaan ServeMux{}
	// mux := http.ServeMux{}

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "Hello, this is the root page.")
	// })

	// mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprint(w, "This is the about page.")
	// })

	// http.ListenAndServe(":8080", &mux)
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	templateUpload := template.Must(template.ParseFiles("templates/uploadImage.html"))
	templateUpload.ExecuteTemplate(w, "uploadImage.html", map[string]string{})
}

func showUploadImage(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("nama")
	gambar, header, _ := r.FormFile("gambar")

	dst, _ := os.Create("./uploads/" + header.Filename)

	io.Copy(dst, gambar)

	templateShowUpload := template.Must(template.ParseFiles("templates/showUploadImage.html"))
	templateShowUpload.ExecuteTemplate(w, "showUploadImage.html", map[string]interface{}{
		"Name":  name,
		"Image": "/static/" + header.Filename,
	})

}

func TestUploadImage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", uploadImage)
	mux.HandleFunc("/upload", showUploadImage)
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./uploads"))))

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}

func downloadImage(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")

	if file == "" {
		fmt.Fprint(w, errors.New("File not found"))
		return
	}

	// Agar langsung terdownload tanpa terrender dibrowser.
	w.Header().Add("Content-Disposition", "attachment; filename=\""+file+"\"")
	http.ServeFile(w, r, "./uploads/"+file)
}

func TestDownloadFile(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", downloadImage)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
