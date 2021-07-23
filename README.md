sysvinit like task runner in go

[![GoDoc](https://godoc.org/github.com/raohwork/gotask?status.svg)](https://godoc.org/github.com/raohwork/gotask)
[![Go Report Card](https://goreportcard.com/badge/github.com/raohwork/gotask)](https://goreportcard.com/report/github.com/raohwork/gotask)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-97%25-brightgreen.svg?longCache=true&style=flat)</a>

# Synopsis

```go
must := func(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
r := gotask.NewRunner()
must(r.QuickAdd("a", initAFunc))
must(r.Add(gotask.New("b", initBFunc, []string{"a"}))

if err := r.RunSync(); err != nil {
    log.Fatal(err)
}
```

### Strange but fairly complex example

In this example, it starts pgsql and redis via systemd, warm-up the cache before start the service.

It uses [routines](https://github.com/raohwork/routines) package

```go
var runner = gotask.NewRunner()

func must(e error) {
    if e != nil {
        log.Fatal(e)
    }
}

func systemdStart(srv string) (err error) {
    if err := exec.Command("systemctl", "start", srv).Run(); err != nil {
        return
    }
    
    // check 3 times every 10 sec
    err = routines.TryAtMost(3, routines.RunAtLeast(
        10 * time.Second, 
        func() error { return exec.Command("systemctl", "status", srv).Run() },
    ))
        
    return
}

func init() {
    must(runner.QuickAdd("pgsql", func() error { return systemdStart("postgresql") }))
    must(runner.QuickAdd("redis", func() error { return systemdStart("redis") }))
}

func fillCache() error {
    // TBD
    return nil
}

func init() {
    must(runner.Add("warm-up", fillCache, []string{
        "prepareConn",
    }))
}

var dbConn *sql.DB
var cacheConn *redis.Client

func prepareConn() error {
    // TBD: setup dbConn and cacheConn
    return nil
}

func init() {
    must(runner.Add("prepareConn", prepareConn, []string{
        "pgsql", "redis",
    }))
}

func main() {
    if err := runner.Run(2); err != nil {
        log.Fatal(runner.Errors())
    }
    
    log.Print(myservicepkg.Start())
}
```

# License

Copyright Chung-Ping Jen <ronmi.ren@gmail.com> 2021-

MPL v2.0
