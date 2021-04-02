package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

type DebugStats struct {
	Time     int64
	Memory   int
	Clients  int
	Profiles map[string]int
}

var debugProfiles = []string{
	"goroutine",
	"heap",
	"block",
	"mutex",
}

var tmplDebug = template.Must(template.New("").Parse(`<!doctype html>
<html>
<head>
	<title>Debug stats</title>
	<meta charset="utf-8">
</head>
<body>
	<table>
		<tr><td>Memory</td><td>{{.Memory}} MB</td></tr>
		<tr><td>Clients</td><td>{{.Clients}}</td></tr>
		{{range $p, $c := .Profiles}}
		<tr><td><a href="/debug/pprof?profile={{$p}}">{{$p}}</a></td><td>{{$c}}</td></tr>
		{{end}}
	</table>
	<br />
	<a href="?raw=1">JSON data</a><br />
	<a href="/debug/profile">CPU Profile</a><br />
	<a href="/debug/pprof?profile=heap&raw=1">Heap Profile</a><br />
	<a href="https://github.com/golang/go/wiki/Performance">Performance Tips</a>
</body>
</html>
`))

func debugAvailable(r *http.Request) bool {
	return strings.HasPrefix(webClientAddr(r), "127.0.0.1")
}

func (s *Server) debugIndex(w http.ResponseWriter, r *http.Request) error {
	if !debugAvailable(r) {
		http.Error(w, "not found", 404)
		return nil
	}
	d := DebugStats{}
	d.Time = time.Now().Unix()
	d.Clients = len(s.getClients())
	d.Profiles = make(map[string]int)
	for _, p := range debugProfiles {
		d.Profiles[p] = pprof.Lookup(p).Count()
	}
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	d.Memory = int(m.Sys / 1000 / 1000)
	if r.FormValue("raw") != "" {
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.SetIndent("", "\t")
		return e.Encode(d)
	}
	return tmplDebug.Execute(w, d)
}

func (s *Server) debugProfile(w http.ResponseWriter, r *http.Request) error {
	if !debugAvailable(r) {
		http.Error(w, "not found", 404)
		return nil
	}
	// Stolen from: https://golang.org/src/net/http/pprof/pprof.go#L116
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="profile"`)
	if err := pprof.StartCPUProfile(w); err != nil {
		return fmt.Errorf("CPU profiling error: %s", err)
	}
	var clientGone <-chan bool
	if cn, ok := w.(http.CloseNotifier); ok {
		clientGone = cn.CloseNotify()
	}
	select {
	case <-time.After(time.Duration(30) * time.Second):
	case <-clientGone:
	}
	pprof.StopCPUProfile()
	return nil
}

func (s *Server) debugPprof(w http.ResponseWriter, r *http.Request) error {
	if !debugAvailable(r) {
		http.Error(w, "not found", 404)
		return nil
	}
	profile := strings.TrimSpace(r.FormValue("profile"))
	found := false
	for _, p := range debugProfiles {
		if profile == p {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid pprof profile: %s", profile)
	}
	p := pprof.Lookup(profile)
	if r.FormValue("raw") != "" {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, profile))
		return p.WriteTo(w, 0)
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		return p.WriteTo(w, 1)
	}
}
