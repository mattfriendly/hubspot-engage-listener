package main

import (
        "crypto/hmac"
        "crypto/sha256"
        "encoding/hex"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "strings"
        "time"

        "github.com/fatih/color"
)

func main() {
        host := os.Getenv("HOST") // Get the host from environment variable
        port := os.Getenv("PORT")
        if port == "" {
                port = "8771" // Default to port 3000
        }

        address := host + ":" + port // Construct the address

        clientSecret := os.Getenv("CLIENT_SECRET")
        allowList := strings.Split(os.Getenv("ALLOW_LIST"), ",")

        mux := http.NewServeMux()
        mux.HandleFunc("/hubspot-webhook", webhookHandler(clientSecret, allowList))

        certPath := "/etc/ssl/linode/fullchain.pem"
        keyPath := "/etc/ssl/linode/privkey.pem"

        log.Printf("Server listening on %s\n", address)
        log.Fatal(http.ListenAndServeTLS(address, certPath, keyPath, mux))
}

func webhookHandler(clientSecret string, allowList []string) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                cyan := color.New(color.FgCyan).SprintfFunc()
                yellow := color.New(color.FgYellow).SprintfFunc()
                green := color.New(color.FgGreen).SprintfFunc()

                log.Println(cyan("Received webhook [%s]: %s %s from %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path, r.RemoteAddr))

                // Log headers
                for name, values := range r.Header {
                        for _, value := range values {
                                log.Println(yellow("Header: %s: %s", name, value))
                        }
                }

                body, err := ioutil.ReadAll(r.Body)
                if err != nil {
                        log.Println("Error reading request body:", err)
                        w.WriteHeader(http.StatusInternalServerError)
                        fmt.Fprintf(w, "Internal Server Error")
                        return
                }

                // Log body
                log.Println(green("Body: %s", string(body)))

                // Verify IP allowlist
                clientIP := strings.Split(r.RemoteAddr, ":")[0]
                if !ipAllowed(clientIP, allowList) {
                        log.Println("IP not allowed:", clientIP)
                        w.WriteHeader(http.StatusForbidden)
                        fmt.Fprintf(w, "IP not allowed")
                        return
                }

                // Verify signature
                signature := r.Header.Get("X-HubSpot-Signature")
                if !verifySignature(signature, clientSecret, r.Method, r.URL.Path, body) {
                        log.Println("Invalid signature")
                        w.WriteHeader(http.StatusForbidden)
                        fmt.Fprintf(w, "Invalid signature")
                        return
                }

                // Process webhook payload
                log.Println("Webhook processed successfully")
                w.WriteHeader(http.StatusOK)
                fmt.Fprintf(w, "Webhook processed successfully")
        }
}

func verifySignature(signature, appSecret, method, path string, body []byte) bool {
        h := hmac.New(sha256.New, []byte(appSecret))
        h.Write([]byte(method + path + string(body)))
        expectedMAC := hex.EncodeToString(h.Sum(nil))
        return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

func ipAllowed(clientIP string, allowList []string) bool {
        for _, ip := range allowList {
                if ip == "0.0.0.0" {
                        return true  // Allow all IPs if "0.0.0.0" is in the list
                }
                if clientIP == ip {
                        return true
                }
        }
        return false
}
