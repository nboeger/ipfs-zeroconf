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

var (
  debug   = flag.Bool("d",false, "Turn debugging on")
  timeout = flag.Int("t", TIMEOUT, "Time out used for making network connections")
  service = flag.String("s", SERVICE, "IPFS Service string to be used")
  domain  = flag.String("D", DOMAIN, "Domain to announce")
)

/*
 * Announce
 * This will announce the IPFS services on the locally 
 * attched network
 */
func Announce(ni NetInfo) {
  for {
    for i, _  := range ni.PublishAddrs {
      id_str := fmt.Sprintf("ID=%s", ni.ID)
      addr_str := fmt.Sprintf("Address=%s", ni.PublishAddrs[i])
      s, err := bonjour.Register("IPFS",
                                 *service,
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
    time.Sleep(time.Duration(*timeout) * time.Second)
  }
}

/*
 * Discover
 * This will look for other IPFS zeroconf servers making
 * announcements. It will then add them to the local 
 * IPFS service
 */
func Discover(ni NetInfo) {
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

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*timeout))
    defer cancel()

    err = resolver.Browse(ctx, *service, *domain, entries)
    if err != nil {
      log.Fatal("Failed to browse network: ", err.Error())
    }
    time.Sleep(time.Duration(*timeout) * time.Second)
    <-ctx.Done()
  }
}


func main() {
  var ni NetInfo

  flag.Parse()

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
  go Discover(ni)
  go Announce(ni)

  select {}
}

