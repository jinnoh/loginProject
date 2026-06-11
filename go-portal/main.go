package main

//no third party dependencies for now just native go libraries
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
//hashmap to map tokens to usernames
var activeSession = make(map[string]string)

//hardcoded credentials for demo
const targetUser = "admin"
const targetPw = "SecPassword123"

//login page
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

//successful login page
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

//unauthroized access page
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
	//go's crypto/rand pacakage generates pseudorandom number and puts it into buffer
	rand.Read(b)
	//encode the random bytes into a url safe base64 string
	return base64.URLEncoding.EncodeToString(b)
}

func main() {
	//get port from environment vairable (GCP cloud run dynmially sets $PORT env var)
	//otherwise default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//captures all requests to root and redirect to login html page
	//uses http status see other redirect status to login page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	//login logic
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		//post req, parse inputs
		if r.Method == http.MethodPost {
			username := r.FormValue("username")
			password := r.FormValue("password")

			//if login valid, we make a session token and store it in the map with username
			if username == targetUser && password == targetPw {
				token := generateSessionToken()
				activeSession[token] = username

				//samesite cokoie with http only flag to defend against xss and csrf
				//issues cookie with token value
				//HTTP only flag prevents client side scirpts form accesing or modifying cookie
				//can preven t against xss where cookies get stolen/hijacked
				//same site cookie attirbute is set to strict to prevent csrf, request must originate
				//from same site or else cookie isn't attached to request
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Path:     "/",
				})
				//upon successful authetnication redirect to dashboard
				http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return
			}
			fmt.Fprintf(w, LoginHTML, "<p style='color: red;'>Invalid credentials.</p>")
			return
		}
		fmt.Fprintf(w, LoginHTML, "")
	})

	//auhtorization for dashboard page, checks for cookie
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}

		//verify cookie with active session map
		username, found := activeSession[cookie.Value]
		if !found {
			http.Redirect(w, r, "/unauthorized", http.StatusSeeOther)
			return
		}

		fmt.Fprintf(w, DashboardHTML, username)
	})

	//unauthroized page handler
	http.HandleFunc("/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, UnauthorizedHTML)
	})

	//deletes session token basically hadnles logout/session ending
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		//delete cookie form map
		if err == nil {
			delete(activeSession, cookie.Value)
		}
		//instruct browser to delte cookie by setting max age to -1
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}