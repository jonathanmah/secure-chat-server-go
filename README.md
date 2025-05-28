# Real-Time Go WebSocket Chat App

**Real-time chat server with secure authentication and WebSocket support**  
Built with Go, PostgreSQL, JavaScript, and Google OAuth 2.0.

<img width="1507" alt="Image" src="https://github.com/user-attachments/assets/67ec2a12-a940-4b6b-b170-dd3081c43c30" />
 
## Features

- Real-time communication via WebSockets with support for multiple chat rooms
- Scalable Pub/Sub architecture using a centralized hub for managing client connections
- Secure HTTP-only cookie sessions with JWT access tokens and refresh tokens stored in Postgres
- Authentication flows included
  - User registration with email verification for activating accounts
  - Login with traditional sign-in or using Google OAuth 2.0
  - Password reset via email  
- Robust error handling and logging for tracking connected clients and messages

## Stack
- Language: Go
- Routing: [chi](https://github.com/go-chi/chi)
- WebSockets: [gorilla/websocket](https://github.com/gorilla/websocket)
- OAuth: [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- Database: PostgreSQL via [lib/pq](https://github.com/lib/pq)
- Frontend: Vanilla JavaScript and WebSocket API

## Getting Started
### 1. Clone the Repo

```bash
git clone https://github.com/jonathanmah/secure-chat-server-go.git
cd secure-chat-server-go
```
### 2. Set up an OAuth 2.0 client
Details on how to create an OAuth 2.0 client and access client id and secret can be found [here](https://support.google.com/googleapi/answer/6158849?hl=en)

### 3. Create a new email
Create an email account to be used for sending automated emails. If you're using Gmail you may need to enable 2-Step Verification and generate an [App Password](https://support.google.com/mail/answer/185833?hl=en) for SMTP to work.

### 4. Create a .env file in project root and add the following environment variables

Create a `.env` file and add the following environment variables:

```env
# SERVER
PORT = 8080 # Port number for HTTP server to listen on
BASE_URL = http://localhost:8080 # Full base URL, used for frontend links and backend routes

# AUTH
ACCESS_TOKEN_SECRET = <your-secret> # used to sign/verify access JWTs
ACTIVATION_TOKEN_SECRET = <your-secret> #  used to sign/verify JWTs for activating accounts or password reset
OAUTH_CLIENT_ID = <your-OAuth-client-id>.apps.googleusercontent.com  #Your Google OAuth 2.0 Client ID
OAUTH_CLIENT_SECRET = <your-OAuth-client-secret>   # Your Google OAuth 2.0 Client secret 
OAUTH_USER_INFO_URL = https://www.googleapis.com/oauth2/v3/userinfo

# EMAIL
SMTP_FROM= example@gmail.com # The email address your app uses to send confirmation/password reset emails
SMTP_HOST=smtp.gmail.com # Host address of Gmail's SMTP server (you can use whichever)
SMTP_PORT=587 # The port used to connect to the SMTP server
SMTP_PASSWORD = <your-smtp-password>

# POSTGRES
PG_USER = <your-postgres-user>
PG_PASSWORD = <your-postgres-password>
PG_HOST = localhost
PG_PORT = 5432
PG_DBNAME = <your-database-name>
PG_SSL_MODE = disable # leave disabled for dev
PG_DRIVER_NAME = postgres # The SQL driver name to use
```

### 5. Set Up PostgreSQL Database
```bash
psql -U <your-postgres-user> -d <your-database-name> -f schema.sql
```
### 6. Run Server

```bash
cd cmd/server
go run main.go
```


