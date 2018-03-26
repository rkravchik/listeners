package listeners

import (
	"net"
	"net/http"
	"time"
)

// TCPKeepAliveListener is a copy of the respective non-exportable struct from net/http.
// It sets TCP keep-alive timeouts on accepted connections.
// It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually go away.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements net.Listener interface.
func (ln TCPKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
