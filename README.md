tcptest
==========

Start A Network Server On Random Local Port (Port of Perl5's [Test::TCP](https://metacpan.org/release/Test-TCP))

[![Build Status](https://travis-ci.org/lestrrat-go/tcptest.png?branch=master)](https://travis-ci.org/lestrrat-go/tcptest)
[![GoDoc](https://godoc.org/github.com/lestrrat-go/tcptest?status.svg)](https://godoc.org/github.com/lestrrat-go/tcptest)

# SYNOPSIS

```go
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

  // Do what you want to do with memcached

  // Then when you're done, you need to kill it
  cmd.Process.Signal(syscall.SIGTERM)

  // And wait
  server.Wait()
```
