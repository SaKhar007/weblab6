package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const dataFile = "data.json"

var mu sync.Mutex

func main() {

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/api", handleAPI)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Open http://localhost:8080/page1.html to start")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func handleAPI(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		mu.Lock()
		data, err := os.ReadFile(dataFile)
		mu.Unlock()

		if err != nil {
			if os.IsNotExist(err) {
				w.Write([]byte("[]"))
				return
			}
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(data)

	} else if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		mu.Lock()
		err = os.WriteFile(dataFile, body, 0644)
		mu.Unlock()

		if err != nil {
			http.Error(w, "Error writing file", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Data saved via Go",
		})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
