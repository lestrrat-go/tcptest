package tcptest

import (
  "fmt"
  "log"
  "net"
  "os/exec"
  "testing"
  "time"
)

func Example() {
  memd := func(port int) {
    cmd := exec.Command("memcached", "-p", fmt.Sprintf("%d", port))
    cmd.Run()
  }

  port, err := Start(memd, 30 * time.Second)
  if err != nil {
    log.Fatalf("Failed to start memcached: %s", err)
  }

  log.Printf("memcached started on port %d", port)
}

func TestBasic(t *testing.T) {
  cb := func(port int) {
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
      t.Fatalf("Failed to listen on port %d: %s", port, err)
    }

    _, err = l.Accept()
    if err != nil {
      t.Fatalf("Failed to accept connection on %d: %s", port, err)
    }
  }

  t.Logf("Starting callback")
  port, err := Start(cb, time.Minute)
  if err != nil {
    log.Fatalf("Failed to start listening on random port: %s", err)
  }

  t.Logf("Attempting to connect to port %d", port)
  _, err = net.Dial("tcp", fmt.Sprintf(":%d", port))
  if err != nil {
    log.Fatalf("Failed to connect to port %d: %s", port, err)
  }

  t.Logf("Successfully connected to port %d", port)
}