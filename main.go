package main

import (
	"flag"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"syscall"
	"time"
  "fmt"
  "context"
  "io/ioutil"
  "github.com/grandcat/zeroconf"
  "github.com/oleksandr/bonjour"
)

const PIDFILE  = "ipfs-zeroconf.pid"
const LOGFILE  = "ipfs-zeroconf.log"
const WORKDIR  = "./"
const PIDPERM  = 0644
const LOGPERM  = 0640
const UMASK    = 027
const ARGS     = "[ipfs-zeroconf]"

const TIMEOUT = 10
const DOMAIN  = "local"
const SERVICE = "_ipfs._tcp"



/*
 * Announce
 * This will announce the IPFS services on the locally 
 * attched network
 */
func Announce(ni NetInfo) {
  log.Print("Announcing timeout: ", *timeout)
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
}

/*
 * Discover
 * This will look for other IPFS zeroconf servers making
 * announcements. It will then add them to the local 
 * IPFS service
 */
func Discover(ni NetInfo) {
  var foo IPFSAddr
  log.Print("Discovering timeout: ", *timeout)
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
      err = ni.AddHost(foo)
      if err != nil {
        // We will print the error but ignore the fail
        // its not a big deal. Just add others
        log.Print(err)
      }
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

var (
  debug   = flag.Bool("d",false, "Turn debugging on")
  timeout = flag.Int("t", TIMEOUT, "Time out used for making network connections")
  service = flag.String("s", SERVICE, "IPFS Service string to be used")
  domain  = flag.String("D", DOMAIN, "Domain to announce")
	signal  = flag.String("cmd", "", `Send command to the daemon: stop — graceful shutdown`)
)

func main() {

	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGQUIT, termHandler)

	cntxt := &daemon.Context{
		PidFileName: PIDFILE,
		PidFilePerm: PIDPERM,
		LogFileName: LOGFILE,
		LogFilePerm: LOGPERM,
		WorkDir:     WORKDIR,
		Umask:       UMASK,
		Args:        flag.Args(),
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := cntxt.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}
		daemon.SendCommands(d)
		return
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatalln(err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Println("- - - - - - - - - - - - - - -")
	log.Println("started")
  go worker(*timeout, *debug)

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}

	log.Println("daemon terminated")
}

var (
	stop = make(chan int)
	done = make(chan int)
)

func worker(timeout int, debug bool) int {
  var ni NetInfo

  err := ni.Setup()
  if err != nil {
    log.Print(err)
    os.Exit(1)
  }

  if debug {
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
  } else {
    log.SetFlags(0)
    log.SetOutput(ioutil.Discard)
  }

	for {
		time.Sleep(time.Second)
		select {
		case <-stop:
	    done <- 1
      return 1
		default:
      go Discover(ni)
      go Announce(ni)
      time.Sleep(time.Duration(timeout) * time.Second)
		}
	}
  return 0
}

func termHandler(sig os.Signal) error {
	log.Println("terminating...")
	stop <- 1
	<-done
	return daemon.ErrStop
}
