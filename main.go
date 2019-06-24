package main

import (
  "log"
  "os"
  "time"
  "flag"
  "fmt"
  "context"
  "github.com/grandcat/zeroconf"
  "github.com/oleksandr/bonjour"
)

const TIMEOUT = 10
const DOMAIN  = "local"
const SERVICE = "_ipfs._tcp"

func Announce(ni NetInfo) {
  for {
    log.Println("Announcing")
    for i, _  := range ni.PublishAddrs {
      id_str := fmt.Sprintf("ID=%s", ni.ID)
      addr_str := fmt.Sprintf("Address=%s", ni.PublishAddrs[i])
      s, err := bonjour.Register("IPFS",
                                 "_ipfs._tcp",
                                 "",
                                 ni.Port[i],
                                 []string{id_str, addr_str},
                                 nil )
      if err != nil {
        log.Fatal(err)
      }
      time.Sleep(1 * time.Second)
      s.Shutdown()
    }
    time.Sleep(TIMEOUT * time.Second)
  }
}

func discover(ni NetInfo) {
  var foo IPFSAddr

  for {
    resolver, err := zeroconf.NewResolver(nil)
    if err != nil {
      log.Fatal("Failed to initialize resolver ", err.Error())
    }
    entries := make(chan *zeroconf.ServiceEntry)

    go func(results <-chan *zeroconf.ServiceEntry) {
      for entry := range results {
        err := foo.ParseBonjour(entry.Text)
        if err != nil {
          log.Fatal("Failed to parse Bonjour string!")
        }
        ni.AddHost(foo)
      }
    }(entries)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT)
    defer cancel()

    err = resolver.Browse(ctx, SERVICE, DOMAIN, entries)
    if err != nil {
      log.Fatal("Failed to browser: ", err.Error())
    }
    time.Sleep(TIMEOUT * time.Second)
    <-ctx.Done()
  }
}

func main() {
  var debug = flag.Bool("d",false, "Turn debugging on") 
  var ni NetInfo
  //var service = flag.String("IPFS", "_ipfs._tcp", "IPFS service")
  //var domain = flag.String("domain", "local", "Look only on local domain")
  //var timeout = flag.Int("wait", 10, "will wait 10s to run discovery")

  flag.Parse()
  log.Print(*debug)
  if *debug {
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
  } else {
    log.SetFlags(0)
  }

  err := ni.Setup()
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }
  go discover(ni)
  go Announce(ni)

  select {}
}

