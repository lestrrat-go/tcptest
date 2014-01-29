package tcptest

import (
  "time"
  "github.com/lestrrat/go-tcputil"
)

type TCPTest struct {
  port int
  exit chan bool
}

func Start(cb func (int), dur time.Duration) (*TCPTest, error) {
  p, err := tcputil.EmptyPort()
  if err != nil {
    return nil, err
  }

  c := make(chan bool, 1)
  go func(c chan bool, p int) {
    cb(p)
    c <-true
  }(c, p)

  err = tcputil.WaitLocalPort(p, dur)
  if err != nil {
    return nil, err
  }

  return &TCPTest { p, c }, nil
}

func (t *TCPTest) Port() int {
  return t.port
}

func (t *TCPTest) Wait() {
  <-t.exit
}
