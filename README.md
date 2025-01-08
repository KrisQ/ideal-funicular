# 🐦 Chirpy: Your Friendly Microblogging Platform

Welcome to **Chirpy**, the lightweight, fast, and charmingly quirky microblogging platform! Whether you're here to chirp, code, or just curious about our chirp-worthy architecture, you're in the right place.

---

## 🚀 Quick Start

### Prerequisites:
- Go 1.23.4 or later
- PostgreSQL
- `.env` file with:
  - `DB_URL`: Your database connection string
  - `SECRET`: Your JWT secret
  - `PLATFORM`: `"dev"` or `"prod"`
  - `POLKA_KEY`: API key for webhooks

### Get Chirping:
1. Clone the repo:  
   ```bash
   git clone https://github.com/yourusername/chirpy.git
   cd chirpy
   ```

2. Install dependencies:  
   ```bash
   go mod tidy
   ```

3. Set up the database:  
   ```bash
   goose up
   ```

4. Run the server:  
   ```bash
   go run main.go
   ```

5. Chirp away at [http://localhost:8080](http://localhost:8080)!

---

## 🛠 Features

- **User Management:** Create, update, or delete users with secure authentication.
- **Microblogging:** Post, retrieve, and delete chirps (max 140 characters, of course).
- **Admin Tools:** Reset the database, view metrics, and manage user upgrades.
- **Dev Mode:** Special endpoints reserved for developers (chirp responsibly!).

---

## 🤖 API Endpoints

Here are some endpoints to get you started:

| Method | Endpoint                  | Description                      |
|--------|---------------------------|----------------------------------|
| GET    | `/api/chirps`             | Retrieve all chirps             |
| POST   | `/api/chirps`             | Create a new chirp              |
| DELETE | `/api/chirps/{chirpId}`   | Delete a chirp                  |
| POST   | `/api/users`              | Create a user                   |
| POST   | `/api/login`              | Log in and get your token       |

---

## 🧪 Running Tests

Because even chirps need quality control:  
```bash
go test ./...
```

---

## 🎨 Code Highlights

- **Middleware:** Handle metrics and developer-only access with ease.
- **SQLC:** Auto-generated database queries keep the codebase tidy.
- **Testing:** MockDB ensures your chirps are always tweet-perfect.

