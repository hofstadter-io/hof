package yagu

import "net"

// Hacky(?) way to ask the system for a free open port that is ready to use.
// call this at the last possible moment before using this port for real.
func GetFreePort() (int, error) {
	// make an address with port 0
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	// listen, thus getting assigned an open port
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	// close when we leave
	defer l.Close()
	
	// return the port we were given
	return l.Addr().(*net.TCPAddr).Port, nil

	// could this port be consumed by the time we use it?
}
