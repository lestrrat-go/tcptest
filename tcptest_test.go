package tcptest

import (
  "fmt"
  "log"
  "net"
  "os/exec"
  "syscall"
  "testing"
  "time"
)

func Example() {
  var cmd *exec.Cmd
  memd := func(port int) {
    cmd = exec.Command("memcached", "-p", fmt.Sprintf("%d", port))
    cmd.SysProcAttr = &syscall.SysProcAttr {
      Setpgid: true,
    }
    cmd.Run()
  }

  server, err := Start(memd, 30 * time.Second)
  if err != nil {
    log.Fatalf("Failed to start memcached: %s", err)
  }

  log.Printf("memcached started on port %d", server.Port())
  defer func() {
    if cmd != nil && cmd.Process != nil {
      cmd.Process.Signal(syscall.SIGTERM)
    }
  }()

  // Do what you want...


  // When you're done, you need to kill the daemon
  // because, well, it doesn't unless you tell it to!
  cmd.Process.Signal(syscall.SIGTERM)

  // Wait for the callback goroutine to exit
  server.Wait()
}

func TestBasic(t *testing.T) {
  cb := func(port int) {
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
      t.Fatalf("Failed to listen on port %d: %s", port, err)
    }

    // We'll be accessed twice
    for i := 0; i < 2; i++ {
      conn, err := l.Accept()
      if err != nil {
        t.Fatalf("Failed to accept connection on %d: %s", port, err)
      }
      time.Sleep(500 * time.Millisecond)
      conn.Close()
    }
    l.Close()
  }

  t.Logf("Starting callback")
  server, err := Start(cb, time.Minute)
  if err != nil {
    log.Fatalf("Failed to start listening on random port: %s", err)
  }

  t.Logf("Attempting to connect to port %d", server.Port())
  conn, err := net.Dial("tcp", fmt.Sprintf(":%d", server.Port()))
  if err != nil {
    log.Fatalf("Failed to connect to port %d: %s", server.Port(), err)
  }
  defer conn.Close()

  t.Logf("Successfully connected to port %d", server.Port())

  server.Wait()
}


func TestMemcached(t *testing.T) {
  // Only run this if we can find a memcached binary in our PATH
  fqname, err := exec.LookPath("memcached")
  if err != nil {
    t.Skip("No memcached available, skipping test")
  }

  t.Logf("Using memcached in %s", fqname)

  var cmd *exec.Cmd
  cb := func(port int) {
    cmd = exec.Command("memcached", "-p", fmt.Sprintf("%d", port))
    cmd.Run()
  }

  server, err := Start(cb, time.Minute)
  if err != nil {
    log.Fatalf("Failed to start listening on random port: %s", err)
  }

  cmd.Process.Signal(syscall.SIGTERM)

  time.Sleep(5 * time.Second)

  _, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", server.Port()))
  if err == nil {
    t.Errorf("After 5 seconds, we can still connect to port %d. Not good", server.Port())
  }
}
