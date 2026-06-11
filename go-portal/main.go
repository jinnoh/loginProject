package main

import(
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/base64"
	"crypto/rand"
)

//store gloablk session in memory
var activeSession = make(map[string]string)

//hardcoded credentials for demo
const targetUser = "admin"
const targetPw = "SecPassword123"

const LoginHTML = `<!DOCTYPE html>
<html>
<head><title>RxDiet Secure Portal - Login</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px;">
    <h2>RxDiet Core Portal Login</h2>
    %s
    <form method="POST" action="/login">
        <label>Username:</label><br>
        <input type="text" name="username" required><br><br>
        <label>Password:</label><br>
        <input type="password" name="password" required><br><br>
        <button type="submit">Login</button>
    </form>
</body>
</html>`

const DashboardHTML = `<!DOCTYPE html>
<html>
<head><title>Dashboard</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px;">
    <h2 style="color: green;">Welcome to the Secure Dashboard!</h2>
    <p>User Identity Status: <strong>Authenticated as %s</strong></p>
    <p>This data is protected behind server-side session tokens.</p>
    <a href="/logout">Logout</a>
</body>
</html>`

const UnauthorizedHTML = `<!DOCTYPE html>
<html>
<head><title>403 - Unauthorized</title></head>
<body style="font-family: Arial, sans-serif; margin: 50px; text-align: center;">
    <h1 style="color: red;">403 - Access Denied</h1>
    <p>You must be logged in to view this page resource.</p>
    <p><a href="/login">Return to Login Page</a></p>
</body>
</html>`

//random session token
func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username := strings.TrimSpace(r.FormValue("username"))
			password := r.FormValue("password")

			if username == targetUser && password == targetPw {
				token := generateSessionToken()
				activeSession[token] = username

				//samesite cokoie with http only flag to defend against xss and csrf
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
			fmt.Fprintf(w, LoginHTML, "<p style='color: red;'>Invalid credentials.</p>")
			return
		}
		fmt.Fprintf(w, LoginHTML, "")
	})

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		//access control cehck for sesion token
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}

		username, found := activeSession[cookie.Value]
		if !found {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}

		fmt.Fprintf(w, DashboardHTML, username)
	})

	http.HandleFunc("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, UnauthorizedHTML)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			delete(activeSession, cookie.Value)
		}
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	log.Printf("Server initializing security boundaries on port %s...", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}