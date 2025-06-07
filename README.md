# Real-Time Go WebSocket Chat App

**Real-time chat server with secure authentication and WebSocket support**  
Built with Go, PostgreSQL, and Google OAuth.
 
<img width="1499" alt="Image" src="https://github.com/user-attachments/assets/9f555e89-ae17-401e-8ad8-4009700069da" />

## Features

- Real-time communication via WebSockets with support for multiple chat rooms
- Scalable pub-sub architecture using a centralized hub for managing client connections
- Secure, HTTP-only cookie-based user sessions with JWT  
- Authentication flows included
  - User registration with email verification for activating accounts
  - Login with traditional sign in or using Google OAuth 2.0
  - Password reset via email  
- Robust error handling and logging for authentication endpoints and tracking connected clients and messages

## Stack
- Language: Go
- Routing: [chi](https://github.com/go-chi/chi)
- WebSockets: [gorilla/websocket](https://github.com/gorilla/websocket)
- OAuth: [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- Database: PostgreSQL via [lib/pq](https://github.com/lib/pq)
- Frontend: vanilla JS and WebSocket API

## Getting Started
### 1. Clone the Repo

```bash
git clone https://github.com/jonathanmah/secure-chat-server-go.git
cd secure-chat-server-go
```
### 2. Setup an OAuth 2.0 client
Details on how to create an OAuth 2.0 client and access client id and secret can be found [here](https://support.google.com/googleapi/answer/6158849?hl=en)

### 3. Create a new email for automation
Create an email account to send automated emails. If using Gmail, you may need to enable 2-Step Verification and generate an [App Password](https://support.google.com/mail/answer/185833?hl=en) for SMTP access.

### 4. Setup environment variables in .env

Create a `.env` file in project root and add the following environment variables:

```env
# SERVER
PORT = 8080
BASE_URL = localhost:8080 # url to app, change this when moving to prod

# AUTH
ACCESS_TOKEN_SECRET = <your-secret> # used to sign JWT
ACTIVATION_TOKEN_SECRET = <your-secret> # used to sign JWT
OAUTH_CLIENT_ID = <your-client-id> # you Google OAuth 2.0 Client ID
OAUTH_CLIENT_SECRET = <your-client-secret> # your Google OAuth 2.0 Client secret 
OAUTH_USER_INFO_URL = https://www.googleapis.com/oauth2/v3/userinfo 

# EMAIL
SMTP_FROM= <your-email-address>
SMTP_HOST=smtp.gmail.com # host address of Gmail SMTP server (you can use whichever)
SMTP_PORT=587
SMTP_PASSWORD = <your-smtp-password>

# POSTGRES
PG_USER = <your-user>
PG_PASSWORD = <your-password>
PG_HOST = <your-host>
PG_PORT = 5432
PG_DBNAME = <your-db-name>
PG_SSL_MODE = disable # SSL mode for PostgreSQL connection
PG_DRIVER_NAME = postgres # The SQL driver name to use to open a connection
```

### 5. Setup Docker and run Docker Compose

```bash
docker-compose up --build
```
