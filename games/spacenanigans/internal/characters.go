package internal

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var tmplChars = template.Must(template.New("chars").Parse(`
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
        <section id="chars">
		<h1>{{.state}}</h1>
{{if eq .state "Characters"}}
		<a href="/chars/create">Create New</a>
		<table>
			<thead>
				<tr><td>Name</td><td>Credits</td><td>Actions</td></tr>
			</thead>
			<tbody>
			{{range .chars}}
				<tr>
					<td>{{.FullName}}</td>
					<td>${{.Credits}}</td>
					<td><a href="/chars/customize?id={{.ID}}">Customize</a> <a href="/chars/delete?id={{.ID}}">Delete</a></td>
				</tr>
			{{end}}
		</table>
		<a href="/">Back</a>
{{else if eq .state "New Character"}}
		<form method="post" action="">
		<input type="text" name="first" value="{{.first}}" placeholder="First name" maxlength="30" required autofocus>
		<input type="text" name="last" value="{{.last}}" placeholder="Last name" maxlength="30" required>
		<select name="gender">
		{{if .gender}}
			<option value="female" {{if eq .gender "female"}}selected{{end}}>Female</option>
			<option value="male" {{if eq .gender "male"}}selected{{end}}>Male</option>
			<option value="other" {{if eq .gender "other"}}selected{{end}}>Other (Androgyne, Bigender, Non-binary, Pangender, Trans, Apache Helicopter, etc.)</option>
			{{else}}
			<option value="female">Female</option>
			<option value="male">Male</option>
			<option value="other">Other (Androgyne, Bigender, Non-binary, Pangender, Trans, Apache Helicopter, etc.)</option>
			{{end}}
		</select>
		<button type="submit" name="action" value="name">Random Name</button>
		<button type="submit" name="action" value="create">Create</button>
		</form>
		<a href="/chars/list">Back</a>
{{else if eq .state "Edit Character"}}
		<form method="post" action="">
		<input type="text" name="first" value="{{.first}}" disabled>
		<input type="text" name="last" value="{{.last}}" disabled>
		<input type="text" name="gender" value="{{.gender}}" disabled>
		<button type="submit">Update</button>
		</form>
		<a href="/chars/list">Back</a>
{{else if eq .state "Delete Character"}}
		<form method="post" actions="">
		<input type="text" value="{{.name}}" disabled>
		<button type="submit">Confirm</button>
		</form>
		<a href="/chars/list">Back</a>
{{else if eq .state "Select Character"}}
		<form method="post" action="">
		<select name="id">
		{{range .chars}}
			<option value="{{.ID}}">{{.FullName}}</option>
		{{end}}
		</select>
		<button type="submit">Play</button>
		</form>
		<a href="/">Back</a>
{{end}}
		<div id="error" {{if not .error}}class="hidden"{{end}}>Error: {{.error}}</div>
        </section>
</body>
</html>
`))

func (s *Server) charsList(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Characters",
		"user":  user,
	}
	chars, err := s.db.ListCharacters(user)
	if err != nil {
		Log("charsList() sql: %s", err)
		ctx["error"] = "failed to get characters"
	}
	ctx["chars"] = chars
	return tmplChars.Execute(w, ctx)
}

func (s *Server) charsCreate(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "New Character",
		"user":  user,
	}
	count, err := s.db.CountCharacters(user)
	if err != nil || count >= 10 {
		if err != nil {
			Log("charsCreate() sql: %s", err)
			ctx["error"] = "Failed to get character list"
		} else {
			ctx["error"] = "You have too many characters"
		}
		return tmplChars.Execute(w, ctx)
	}
	// TODO: validate these values
	first := r.FormValue("first")
	last := r.FormValue("last")
	gender := r.FormValue("gender")
	if r.Method != "POST" || r.FormValue("action") == "name" {
		ctx["first"] = RandomFirstName(gender)
		ctx["last"] = RandomLastName()
		ctx["gender"] = gender
		return tmplChars.Execute(w, ctx)
	}
	c := Character{
		Username:  user,
		FirstName: first,
		LastName:  last,
		Birthday:  time.Now(),
		IsMale:    gender == "male",
	}
	err = s.db.CreateCharacter(c)
	if err != nil {
		Log("charsCreate() sql: %s", err)
		ctx["error"] = err
		return tmplChars.Execute(w, ctx)
	}
	http.Redirect(w, r, "/chars/list", http.StatusSeeOther)
	return nil
}

func (s *Server) charsCustomize(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Edit Character",
		"user":  user,
	}
	return tmplChars.Execute(w, ctx)
}

func (s *Server) charsDelete(w http.ResponseWriter, r *http.Request) error {
	user := getCookie(r)
	ctx := map[string]interface{}{
		"state": "Delete Character",
		"user":  user,
	}
	if r.Method != "POST" {
		c, err := s.db.GetCharacter(r.FormValue("id"), user)
		if err != nil {
			Log("charsDelete() sql: %s", err)
			ctx["error"] = "Failed to get character"
		}
		ctx["name"] = c.FullName()
		return tmplChars.Execute(w, ctx)
	}
	err := s.db.DelCharacter(r.FormValue("id"), user)
	if err != nil {
		Log("charsDelete() sql: %s", err)
		ctx["error"] = "Failed to delete character"
		return tmplChars.Execute(w, ctx)
		// TODO: probs should get char and set the name too, like above
	}
	http.Redirect(w, r, "/chars/list", http.StatusSeeOther)
	return nil
}

func (s *Server) charsPlay(w http.ResponseWriter, r *http.Request) error {
	sendTemplate := func() {
		w.Header().Set("Cache-Control", "no-cache")
		s.serveStatic(w, r, "/static/game.html")
	}

	user := getCookie(r)
	c := s.getClient(user)
	if c != nil {
		if c.IsConnected() {
			return fmt.Errorf("client already connected")
		}
		sendTemplate()
		return nil
	}

	ctx := map[string]interface{}{
		"state": "Select Character",
		"user":  user,
	}
	if r.Method != "POST" {
		chars, err := s.db.ListCharacters(user)
		if err != nil {
			Log("charsList() sql: %s", err)
			ctx["error"] = "failed to get characters"
		}
		ctx["chars"] = chars
		return tmplChars.Execute(w, ctx)
	}

	char, err := s.db.GetCharacter(r.FormValue("id"), user)
	if err != nil {
		Log("charsList() sql: %s", err)
		ctx["error"] = "failed to get character"
		return tmplChars.Execute(w, ctx)
		// TODO: same here, probs should get full list of chars like above
	}
	isNew := s.newClient(&char)
	if !isNew {
		return fmt.Errorf("client already connected")
	}
	sendTemplate()
	return nil
}
