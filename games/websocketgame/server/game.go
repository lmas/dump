package server

func StartGame(addr string) {
	cl := NewClientHandler()

	StartWeb(addr, cl.OnRequest)
}
