# URL-Shortener
**A simple, fast, and lightweight URL shortening service written in Go.**

This project provides a robust and easy-to-deploy URL shortening service. It's designed for efficiency and simplicity, making it an excellent starting point for anyone looking to build a microservice with Go or to quickly deploy a self-hosted URL shortener.

## ✨ Features
- **Fast & Efficient:** Leveraging Go's concurrency model and standard library, the service handles requests with minimal overhead.
- **RESTful API:** Provides a clean and predictable REST API for creating and retrieving shortened URLs.
- **SQLite Storage:** All shortened links are saved in a SQLite database, ensuring data persistence.
- **Dockerized:** A `Dockerfile` is included to easily build and containerize the service for simple and portable deployment.
- **Simple Redirects:** Automatically redirects users from the shortened URL to the original long URL.
- **Bruno Collection:** Includes a Bruno collection with examples of all API calls for testing the service (link creation and redirection).

## 🚀 Getting Started
### Prerequisites
- Go (version **1.24.4** or newer)
- Docker (if you intend to use the `Dockerfile`)

### Installation & Run
1. **Clone the repository:**
```shell
git clone https://github.com/The-EpaG/URL-Shortener.git
cd URL-Shortener
```

2. **Install the SQLite dependency:**
```shell
go get github.com/mattn/go-sqlite3
```

3. **Run the application:**
```shell
go run main.go
```

The service will start on port `8080` by default.

## Docker
1. **Build the Docker image:**
```shell
docker build -t url-shortener .
```

2. **Run the container:**
```shell
docker run -p 8080:8080 url-shortener
```

The service will be accessible at http://localhost:8080.

## 📁 Project Structure
```shell
.
├── main.go               # Main application entry point
├── bruno/                # Bruno collection for API testing
├── internal/             # Internal packages for business logic
│   ├── api/              # Handlers for the REST API
│   └── storage/          # Storage interface and SQLite implementation
├── links.db              # The SQLite database containing the links
├── Dockerfile            # File to containerize the service
├── go.mod                # Go module file
├── go.sum                # Checksums for module dependencies
└── README.md             # This file
```

## 🤝 Contributing
Contributions are welcome! If you have any ideas for new features, bug fixes, or improvements, feel free to open an issue or submit a pull request.

## 📄 License
This project is released under the *GPL-3.0* License. For details, see the LICENSE file.