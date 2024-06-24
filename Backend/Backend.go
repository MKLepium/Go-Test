package Backend

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	// read html file
	file, err := os.Open("./static/index.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func postTest(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /api/test request\n")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // I am not sure if this is necessary

	fmt.Printf("Body: %s\n", body)
	// json response
	io.WriteString(w, `{"status": "Request received"}`)
}

func HTTPServer(port int, hostname string) *http.Server {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/api/test", postTest)

	var addr = fmt.Sprintf("%s:%d", hostname, port)
	fmt.Printf("Starting server on %s\n", addr)
	srv := &http.Server{
		Addr: addr,
	}
	return srv
}
