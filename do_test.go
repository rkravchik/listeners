package listeners

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoListenerOnFirstAcceptDo(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	const (
		once = "once"
		regular = "regular"
	)

	address := ":0"
	addr, err := net.ResolveTCPAddr("tcp", address)
	require.NoError(err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(err)

	firstc := make(chan string, 1)
	fadl := NewDoListener(l).OnFirstAcceptDo(func() {
		firstc <- once
	})

	handler := func(w http.ResponseWriter, r *http.Request) {
		select {
		case first := <-firstc:
			fmt.Fprintf(w, first)
		default:
			fmt.Fprint(w, regular)
		}
	}
	ts := &httptest.Server{
		Listener: fadl,
		Config:   &http.Server{Handler: http.HandlerFunc(handler)},
	}
	ts.Start()

	tc := ts.Client()

	resp, err := tc.Get(ts.URL)
	require.NoError(err)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(err)
	resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.EqualValues(once, body)

	resp, err = tc.Get(ts.URL)
	require.NoError(err)
	body, err = ioutil.ReadAll(resp.Body)
	require.NoError(err)
	resp.Body.Close()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.EqualValues(regular, body)
}
