package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Redis client
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// Session struct
type Session struct {
	User string `json:"user"`
	Role string `json:"role"`
}

func main() {
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/flag", getFlag)
	http.HandleFunc("/flag-alt", flagAlt)
	http.HandleFunc("/admin/delete", adminDelete)


	fmt.Println("🚩 RedFlag running on http://localhost:8080")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

// Renders login form
func loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
		<html>
		<head>
			<title>RedFlag Login</title>
			<style>
				body {
					background-color: #0f0f0f;
					color: #00ff99;
					font-family: "Courier New", monospace;
					text-align: center;
					padding-top: 50px;
				}
				img {
					width: 220px;
					margin-bottom: 20px;
					filter: drop-shadow(0 0 10px #ff0000);
				}
				input, button {
					background: #111;
					color: #00ff99;
					border: 1px solid #00ff99;
					padding: 10px;
					font-size: 16px;
				}
				h1 {
					font-size: 26px;
					color: red;
				}
			</style>
		</head>
		<body>
			<img src="/static/redflag.png" alt="RedFlag Logo"/>
			<h1>💀 REDFLAG</h1>
			<p>The Damn Vulnerable Redis Lab</p>
			<form action="/login" method="POST">
				<input name="username" placeholder="username" />
				<button type="submit">Login</button>
			</form>
		</body>
		</html>
	`))
}



// Handles login and creates session in Redis
func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")

	session := Session{User: username, Role: "user"}
	data, _ := json.Marshal(session)

	// Store session with no TTL
	rdb.Set(ctx, "session:"+username, data, 0)

	http.Redirect(w, r, "/dashboard?user="+username, http.StatusFound)
}

// Displays dashboard and flag link
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	val, err := rdb.Get(ctx, "session:"+user).Result()
	if err != nil {
		http.Error(w, "Session not found", 403)
		return
	}

	var s Session
	json.Unmarshal([]byte(val), &s)

	msg := fmt.Sprintf("<h1>Welcome %s</h1>", s.User)
	msg += `<p><a href="/flag?user=` + s.User + `">Get Flag</a></p>`
	msg += `<p><a href="/flag-alt?user=admin">Try to Steal Admin's Flag</a></p>`

	w.Write([]byte(msg))
}

// Secure flag endpoint (requires role = admin)
func getFlag(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	val, err := rdb.Get(ctx, "session:"+user).Result()
	if err != nil {
		http.Error(w, "Session not found", 403)
		return
	}

	var s Session
	json.Unmarshal([]byte(val), &s)

	if s.Role == "admin" {
		flag, err := rdb.Get(ctx, "flag:"+s.User).Result()
		if err != nil {
			http.Error(w, "Flag not found", 404)
			return
		}
		w.Write([]byte("🎉 FLAG: " + flag))
	} else {
		w.Write([]byte("🚫 Access Denied: You do not have permission to view the flag."))
	}
}

// Vulnerable IDOR flag endpoint (no auth)
func flagAlt(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	flag, err := rdb.Get(ctx, "flag:"+user).Result()
	if err != nil {
		http.Error(w, "Flag not found or user invalid", 404)
		return
	}

	msg := fmt.Sprintf("🕵️ You've found %s's flag: %s", user, flag)
	w.Write([]byte(msg)) 
}

func adminDelete(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")
	if user == "" {
		http.Error(w, "Missing ?user= paramater", 400)
		return
	}

	rdb.Del(ctx, "session:"+user)
	rdb.Del(ctx, "flag:"+user)

	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	msg := fmt.Sprintf(`
		<h2 style= "color: red;">⚠️ Target '%s' wiped from memory</h2>
		<p>Session and flag deleted.</p>
		<a href="/dashboard?user=admin">Return to Dashboard</a>
		`, user)
		w.Write([]byte(msg))
	
}
