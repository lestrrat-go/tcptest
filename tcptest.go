package tcptest

import (
  "time"
  "github.com/lestrrat/go-tcputil"
)

func Start(cb func (int), dur time.Duration) (int, error) {
  p, err := tcputil.EmptyPort()
  if err != nil {
    return 0, err
  }

  go cb(p)

  err = tcputil.WaitLocalPort(p, dur)
  if err != nil {
    return 0, err
  }

  return p, nil
}
