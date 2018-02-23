package tcptest

import (
	"time"

	"github.com/lestrrat-go/tcputil"
)

type TCPTest struct {
	port int
	exit chan struct{}
}

func (t *TCPTest) start(cb func(*TCPTest)) {
	defer t.Stop()
	cb(t)
}

func (t *TCPTest) Done() <-chan struct{} {
	return t.exit
}

// Stop sends a notification through t.Done() that the server should be
// stopped. It's up to the provider of the callback function to properly
// exit and clean the server
func (t *TCPTest) Stop() {
	t.exit<-struct{}{}
}

func Start2(cb func(*TCPTest), dur time.Duration) (*TCPTest, error) {
	p, err := tcputil.EmptyPort()
	if err != nil {
		return nil, err
	}

	t := &TCPTest{
		port: p,
		exit: make(chan struct{}, 10), // random buffer number
	}
	go t.start(cb)

	if err := tcputil.WaitLocalPort(p, dur); err != nil {
		close(t.exit)
		return nil, err
	}

	return t, nil
}

// Start is left for backward compatibility, but you should use
// Start2(), which allows better management of the lifecycle
func Start(cb func(int), dur time.Duration) (*TCPTest, error) {
	p, err := tcputil.EmptyPort()
	if err != nil {
		return nil, err
	}

	c := make(chan struct{}, 10)
	go func(c chan struct{}, p int) {
		defer func() { c <- struct{}{} }()
		cb(p)
	}(c, p)

	err = tcputil.WaitLocalPort(p, dur)
	if err != nil {
		return nil, err
	}

	return &TCPTest{p, c}, nil
}

func (t *TCPTest) Port() int {
	return t.port
}

func (t *TCPTest) Wait() {
	<-t.exit
}
