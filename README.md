1. didalam golang web sudah include web server jadi tidak membutuhkan web server lagi seperti apache atau nginx.

2. 

3. handler 
    - digunakan untuk menangani request dari client
    - tetapi sayangnya handler ini hanya dapat menangani 1 request saja.

4. serveMux 
    - sama seperti handler, serverMux ini juga digunakan untuk menangani request dari client.
    - bedanya serveMux dapat menangani banyak request url.
    - dalam bahasa lain.. ini sama seperti router.

9. Header
    - digunakan untuk membuat sesuatu yang dibutuhkan untuk memasuki website, bisa terdiri dari auth, token, dan sebagainya. ada setandar penamaan di header