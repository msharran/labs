package main

func main() {
	svr := NewTCPServer(TCPServerConfig{
		Host: "localhost",
		Port: "7993",
	})
	svr.ListenAndServe()
}
