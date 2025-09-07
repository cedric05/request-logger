package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize SQLite DB
	db, err := sql.Open("sqlite3", "proxy_logs.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	createTable := `CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		request_received_ts TEXT,
		request_sent_to_proxy_ts TEXT,
		response_received_from_proxy_ts TEXT,
		response_sent_to_client_ts TEXT,
		response_end_ts TEXT,
		url TEXT,
		method TEXT,
		req_headers TEXT,
		req_body TEXT,
		resp_headers TEXT,
		resp_body TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Director modifies the request before sending to backend
		},
		ModifyResponse: func(resp *http.Response) error {
			return nil
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		url := r.URL.String()
		method := r.Method
		reqHeaders := r.Header
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		reqBody := string(bodyBytes)
		requestReceivedTs := start.Format(time.RFC3339Nano)

		// Log request received
		requestSentToProxyTs := time.Now().Format(time.RFC3339Nano)

		// Proxy the request
		proxy.ServeHTTP(w, r)

		responseSentToClientTs := time.Now().Format(time.RFC3339Nano)
		responseEndTs := responseSentToClientTs // For now, same as sent to client

		// For demonstration, response headers/body are not captured here
		_, err := db.Exec(`INSERT INTO logs (
			request_received_ts, request_sent_to_proxy_ts, response_received_from_proxy_ts, response_sent_to_client_ts, response_end_ts, url, method, req_headers, req_body, resp_headers, resp_body
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			requestReceivedTs, requestSentToProxyTs, "", responseSentToClientTs, responseEndTs, url, method, reqHeaders, reqBody, "", "")
		if err != nil {
			log.Printf("Failed to log request: %v", err)
		}
	})

	log.Println("Proxy server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
