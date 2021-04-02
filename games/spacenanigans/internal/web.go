package internal

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gobwas/ws"
)

var (
	gTimeout  = time.Duration(60) * time.Second
	gUpgrader = &ws.HTTPUpgrader{
		Timeout: gTimeout,
	}
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	switch r.URL.Path {
	case "/":
		h = s.Handler(s.pageIndex)

	case "/auth/register":
		h = s.mwLoginNotRequired(s.Handler(s.authRegister))
	case "/auth/login":
		h = s.mwLoginNotRequired(s.Handler(s.authLogin))
	case "/auth/logout":
		h = s.mwLoginRequired(s.Handler(s.authLogout))
	case "/auth/account":
		h = s.mwLoginRequired(s.Handler(s.authAccount))
	case "/auth/delete":
		h = s.mwLoginRequired(s.Handler(s.authDelete))

	case "/chars/list":
		h = s.mwLoginRequired(s.Handler(s.charsList))
	case "/chars/create":
		h = s.mwLoginRequired(s.Handler(s.charsCreate))
	case "/chars/customize":
		h = s.mwLoginRequired(s.Handler(s.charsCustomize))
	case "/chars/delete":
		h = s.mwLoginRequired(s.Handler(s.charsDelete))
	case "/chars/play":
		h = s.mwLoginRequired(s.Handler(s.charsPlay))

	case "/game/ws":
		h = s.mwLoginRequired(s.Handler(s.connectClient))

	case "/debug":
		h = s.Handler(s.debugIndex)
	case "/debug/profile":
		h = s.Handler(s.debugProfile)
	case "/debug/pprof":
		h = s.Handler(s.debugPprof)

	default:
		h = s.StaticHandler(r.URL.Path)
	}
	h.ServeHTTP(w, r)
	r.Body.Close()
}

////////////////////////////////////////////////////////////////////////////////

var tmplIndex = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Spacenanigans</title>
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
	<section id="index">
		<h1>Spacenanigans</h1>
	{{if not .user}}
		<a href="/auth/login">Login</a>
		<a href="/auth/register">Register</a>
	{{else}}
		<p>Welcome {{.user}}</p>
	{{end}}
	</section>
</body>
</html>
`))

func (s *Server) pageIndex(w http.ResponseWriter, r *http.Request) error {
	ctx := map[string]interface{}{
		"user": getCookie(r),
	}
	return tmplIndex.Execute(w, ctx)
}
