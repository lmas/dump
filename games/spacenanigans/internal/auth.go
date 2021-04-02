package internal

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) mwLoginRequired(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getCookie(r)
		if user == "" {
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (s *Server) mwLoginNotRequired(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := getCookie(r)
		if user != "" {
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

////////////////////////////////////////////////////////////////////////////////

var tmplAuth = template.Must(template.New("auth").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{.state}}</title>
	<meta charset="UTF-8">
	<link rel="stylesheet" href="/static/style.css">
</head>
<body>
{{if .user}}
	<header>
		<span><a href="/">Spacenanigans</a></span>
		<a href="/chars/play">Play</a>
		<a href="/chars/list">Characters</a>
		<a href="/auth/account">Account</a>
		<a href="/auth/logout">Logout</a>
	</header>
{{end}}
	<section id="auth">
		<h1>{{.state}}</h1>
		<form method="post" action="">
{{if eq .state "Register"}}
		<input type="text" name="user" placeholder="Username" maxlength="50" required autofocus>
		<input type="password" name="pass" placeholder="Password" maxlength="50" required>
		<input type="text" name="email" placeholder="Email (optional)" maxlength="50">
		<button type="submit">Register</button>
		<a href="/">Back</a>
{{else if eq .state "Login"}}
		<p>Oh noes, we're gonna have to use a cookie to keep you logged in!</p>
		<input type="text" name="user" placeholder="Username" maxlength="50" required autofocus value="alex">
		<input type="password" name="pass" placeholder="Password" maxlength="50" required value="alex">
		<button type="submit">Login</button>
		<a href="/">Back</a>
{{else if eq .state "Logout"}}
		<button type="submit">Confirm</button>
		<a href="/">Back</a>
{{else if eq .state "Account"}}
		<input type="text" value="{{.user}}" disabled>
		<input type="text" name="email" placeholder="Email (optional)" maxlength="50" value="{{.email}}">
		<input type="password" name="pass" placeholder="New Password (optional)" maxlength="100">
		<button type="submit">Update</button>
		<a href="/auth/delete">Delete account</a>
		<a href="/">Back</a>
{{else if eq .state "Delete"}}
		<button type="submit">Confirm</button>
		<a href="/auth/account">Back</a>
{{end}}
		</form>
		<div id="error" {{if not .error}}class="hidden"{{end}}>Error: {{.error}}</div>
	</section>
</body>
</html>
`))

func (s *Server) authRegister(w http.ResponseWriter, r *http.Request) error {
	ctx := map[string]interface{}{
		"state": "Register",
	}
	if r.Method != "POST" {
		return tmplAuth.Execute(w, ctx)
	}
	user := r.FormValue("user")
	pass := r.FormValue("pass")
	if user == "" || pass == "" || len(user) > 50 {
		ctx["error"] = "invalid username or password"
		return tmplAuth.Execute(w, ctx)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		Log("authRegister() bcrypt: %s", err)
		ctx["error"] = "failed to generate password hash"
		return tmplAuth.Execute(w, ctx)
	}
	a := Account{
		Username: user,
		Hash:     string(hash),
		Email:    r.FormValue("email"),
		Created:  time.Now(),
	}
	err = s.db.CreateAccount(a)
	if err != nil {
		Log("authRegister() sql: %s", err)
		ctx["error"] = "failed to save account"
		return tmplAuth.Execute(w, ctx)
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	return nil
}

func (s *Server) authLogin(w http.ResponseWriter, r *http.Request) error {
	ctx := map[string]interface{}{
		"state": "Login",
	}
	if r.Method != "POST" {
		return tmplAuth.Execute(w, ctx)
	}
	user := r.FormValue("user")
	a, _ := s.db.GetAccount(user) // Yap, don't care about any errors here
	err := bcrypt.CompareHashAndPassword([]byte(a.Hash), []byte(r.FormValue("pass")))
	if err != nil {
		ctx["error"] = "invalid username or password"
		return tmplAuth.Execute(w, ctx)
	}
	a.LastLogin = time.Now()
	err = s.db.UpdateAccount(a)
	if err != nil {
		Log("authLogin() sql: %s", err)
	}
	setCookie(w, a.Username)
	w.Header().Set("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (s *Server) authLogout(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Logout",
		"user":  user,
	}
	if r.Method == "GET" {
		return tmplAuth.Execute(w, ctx)
	}
	setCookie(w, "")
	w.Header().Set("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (s *Server) authAccount(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Account",
		"user":  user,
	}
	a, err := s.db.GetAccount(user)
	if err != nil {
		Log("authAccount() sql: %s", err)
		ctx["error"] = "invalid username or password"
		return tmplAuth.Execute(w, ctx)
	}
	if r.Method != "POST" {
		ctx["user"] = a.Username
		ctx["email"] = a.Email
		return tmplAuth.Execute(w, ctx)
	}
	a.Email = r.FormValue("email")
	pass := r.FormValue("pass")
	if pass != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			Log("authAccount() bcrypt: %s", err)
			ctx["error"] = "failed to generate password hash"
			return tmplAuth.Execute(w, ctx)
		}
		a.Hash = string(hash)
	}
	err = s.db.UpdateAccount(a)
	if err != nil {
		Log("authAccount() sql: %s", err)
		ctx["error"] = "failed to save account"
	}
	ctx["user"] = a.Username
	ctx["email"] = a.Email
	return tmplAuth.Execute(w, ctx)
}

func (s *Server) authDelete(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Delete",
		"user":  user,
	}
	if r.Method != "POST" {
		return tmplAuth.Execute(w, ctx)
	}
	err := s.db.DelAccount(user)
	if err != nil {
		Log("authDelete() sql: %s", err)
		ctx["error"] = "failed to delete account"
		return tmplAuth.Execute(w, ctx)
	}
	setCookie(w, "")
	w.Header().Set("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
