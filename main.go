package main

func main() {
	svr := NewTCPServer(TCPServerConfig{
		Host: "localhost",
		Port: "8080",
	})
	svr.ListenAndServe()
}
