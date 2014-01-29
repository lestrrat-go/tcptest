go-tcptest
==========

Start A Network Server On Random Local Port (Port of Perl5's TCP::Test)

```go
  memd := func(port int) {
    cmd := exec.Command("memcached", "-p", fmt.Sprintf("%d", port))
    cmd.Run()
  }

  port, err := Start(memd, 30 * time.Second)
  if err != nil {
    log.Fatalf("Failed to start memcached: %s", err)
  }

  log.Printf("memcached started on port %d", port)
```

API docs: http://godoc.org/github.com/lestrrat/go-tcptest
