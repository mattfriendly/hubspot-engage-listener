# HubSpot Microserver

This repository contains two implementations of a microserver designed to process webhooks from HubSpot with signature verification and IP allow-listing. The server is implemented in **Go** and **Node.js**, providing flexibility for different environments and developer preferences.

---

## Features

- **Webhook Handling:** Listens for webhooks on the `/hubspot-webhook` endpoint.
- **Signature Verification:** Validates webhooks using HMAC SHA-256 signatures.
- **IP Allow-Listing:** Restricts access to the server based on a configurable list of allowed IPs.
- **TLS Support:** Go implementation supports secure communication with TLS certificates.
- **Customizable:** Environment variables allow for easy configuration.

---

## Prerequisites

- **Environment Variables:**
  - `HOST`: The hostname or IP address for the server.
  - `PORT`: The port number for the server (default: 8771 for Go, 3000 for Node.js).
  - `CLIENT_SECRET`: The HubSpot client secret used for signature verification.
  - `ALLOW_LIST`: Comma-separated list of IP addresses allowed to access the server.
- **TLS Certificates (Go implementation):**
  - Place the certificate (`fullchain.pem`) and key (`privkey.pem`) at `/etc/ssl/linode/`.

---

## Setup

### Go Implementation

1. **Install Go:** Ensure Go is installed on your machine.
2. **Clone Repository:**
   ```bash
   git clone <repository-url>
   cd hubspot-microserver-go
   ```
3. Set Environment Variables: Use a .env file or export them directly:   
   ```bash
   export HOST=0.0.0.0
   export PORT=8771
   export CLIENT_SECRET=your_client_secret
   export ALLOW_LIST=127.0.0.1,192.168.1.1
   ```
4. Run the Server:
   ```bash
   go run main.go
   ```
