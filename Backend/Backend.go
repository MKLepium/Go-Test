package Backend

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type APIResponse struct {
	Data   string
	Error  string
	Status string // Could be "pending", "completed", "error", etc.
}

type UserResponse struct {
	Responses map[string]APIResponse
	Status    string
}

var (
	responseStore = make(map[string]*UserResponse)
	mu            sync.Mutex
)

func initializeResponse(id string) {
	mu.Lock()
	defer mu.Unlock()
	responseStore[id] = &UserResponse{
		Responses: make(map[string]APIResponse),
		Status:    "processing",
	}
}

func generateID() string {
	return fmt.Sprintf("%d", rand.Intn(1000000))
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	// read html file
	file, err := os.Open("./Frontend/index.html")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}

func callAPI(url string, id string, key string, requestBody string) {
	// so I can later add the body
	fmt.Printf("Making API call to %s, with Body %s\n", url, requestBody)
	// Simulated API call
	resp, err := http.Get(url)
	if err != nil {
		updateResponse(id, key, APIResponse{Error: err.Error(), Status: "error"})
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		updateResponse(id, key, APIResponse{Error: err.Error(), Status: "error"})
		return
	}

	updateResponse(id, key, APIResponse{Data: string(data), Status: "completed"})
}

func handleUserRequest(w http.ResponseWriter, r *http.Request) {
	id := generateID() // Ensure you have a function to generate unique IDs
	initializeResponse(id)

	// read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// For now two random API calls and a python script to generate a random image.
	go callAPI("https://some-random-api.com/facts/koala", id, "koala_fact", string(body))
	go callAPI("https://catfact.ninja/fact", id, "cat_fact", string(body))
	go callPythonScript(string(body), "./ImageGenerator/ImageGenerator.py", id, "imagepath")

	fmt.Fprintf(w, `{"id": "%s", "message": "Request is being processed"}`, id)
}

func callPythonScript(request string, pythonScriptPath string) (string, error) {
	fmt.Printf("Calling Python script with request: %s\n", request)
	fmt.Printf("Python script path: %s\n", pythonScriptPath)

	// Call the Python script
	cmd := exec.Command("python", pythonScriptPath, request)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error calling Python script: %s\n", err)
		fmt.Printf("Script output: %s\n", string(output))
		return "", err
	}

	// Split the output into lines
	lines := strings.Split(string(output), "\n")
	for i := 1; i < len(lines)-1; i++ {
		// Check if the current line is surrounded by the known marker lines
		if strings.Contains(lines[i-1], "I need something to find the output") &&
			strings.Contains(lines[i+1], "I need something to find the output") {
			uuid := strings.TrimSpace(lines[i])
			fmt.Printf("Filtered Python script output: %s\n", uuid)
			return uuid, nil
		}
	}

	return "", fmt.Errorf("no UUID found in the script output")
}

func handleImageGeneration(request string, pythonScriptPath string, id string, key string) {
	uuid, err := callPythonScript(request, pythonScriptPath)
	if err != nil {
		fmt.Printf("Error Image Generation: %s\n", err)
	}
	mu.Lock()
	responseStore[id].Responses[key] = APIResponse{Data: uuid, Status: "completed"}
	mu.Unlock()
}

func getResponse(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	response, exists := responseStore[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func updateResponse(id string, key string, response APIResponse) {
	mu.Lock()
	defer mu.Unlock()
	if userResp, exists := responseStore[id]; exists {
		userResp.Responses[key] = response
		// Check if all responses are completed
		allCompleted := true
		for _, resp := range userResp.Responses {
			if resp.Status != "completed" {
				allCompleted = false
				break
			}
		}
		if allCompleted {
			userResp.Status = "completed"
		}
	}
}

func HTTPServer(port int, hostname string) *http.Server {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/api/request", handleUserRequest)
	http.HandleFunc("/api/response", getResponse)

	var addr = fmt.Sprintf("%s:%d", hostname, port)
	fmt.Printf("Starting server on %s\n", addr)
	srv := &http.Server{
		Addr: addr,
	}
	return srv
}
