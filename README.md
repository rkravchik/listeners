# Net Listeners

## TCPKeepAliveListener

TCPKeepAliveListener is a copy of the non-exportable respective struct from net/http.

Sometimes it is needed to check port to listen not in `http.ListenAndServe(address)`.

Of course one can do something like this:

    l, err := net.Listen("tcp", addr)
    if err != nil {
        // ...
    }
    // some stuff
    // ...
    http.Serve(l)

But if we look at http.Server.ListenAndServe() a little closer we will see:

    ln, err := net.Listen("tcp", addr)
	if err != nil {
        return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})

And that's why if you want:

1. try to listen port before http.Serve;
2. keep-alive timeouts on accepted connections;

you have to use listeners like `TCPKeepAliveListener`.

Simple usage example:

    package main
     
    import "github.com/rkravchik/listeners"
     
    func main() {
        // ...
        addr, err := net.ResolveTCPAddr("tcp", address)
        if err != nil {
            // ...
        }
        l, err := net.ListenTCP("tcp", addr)
        if err != nil {
            // ...
        }
        err = http.Serve(listeners.TCPKeepAliveListener{TCPListener: l})
        // ...
    }
