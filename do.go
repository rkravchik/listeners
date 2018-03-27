package listeners

import (
	"net"
	"sync"
	"time"
)

// DoListener is a net.TCPListener with keep-alive timeouts and
// ability to execute the given func on first accept being called.
type DoListener struct {
	*net.TCPListener
	once sync.Once
	do   func()
	fn   func() (net.Conn, error)
}

// NewDoListener instance.
func NewDoListener(l *net.TCPListener) *DoListener {
	ln := &DoListener{TCPListener: l}
	ln.fn = ln.firstAccept
	return ln
}

// Accept implements net.Listener interface.
func (ln *DoListener) Accept() (net.Conn, error) {
	return ln.fn()
}

// OnFirstAcceptDo sets func that will be fired on the first Accept execution.
// Use it like:
//	l, err := net.ListenTCP("tcp", addr)
//	// catch err
//	dl := NewDoListener(l).OnFirstAcceptDo(func() {fmt.Println("i'm ready to accept")})
func (ln *DoListener) OnFirstAcceptDo(fn func()) *DoListener {
	ln.do = fn
	return ln
}

// firstAccept atomically changes accept method and executes Do func.
func (ln *DoListener) firstAccept() (net.Conn, error) {
	ln.once.Do(func() {
		ln.fn = ln.accept
		ln.do()
	})
	return ln.fn()
}

// accept internal method that implements net.Listener interface.
func (ln *DoListener) accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
