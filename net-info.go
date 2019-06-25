package main

import (
 "log"
  "net"
  "context"
  "time"
  shell "github.com/ipfs/go-ipfs-api"
)

type NetInfo struct {
  ID           string
  IPS          []net.IPNet
  Port         []int
  PublishAddrs []string
}

/*
 * Setup 
 * This will populate the struct with
 * the local IP's and exclude the
 * localhost
 */
func (ni *NetInfo) Setup() error {
  if len(ni.PublishAddrs) != 0 {
    return nil
  }

  // first get local addresses
  err := ni.my_addresses()
  // now populate with IPFS addresses
  sh := shell.NewShell("localhost:5001")
  foo, err := sh.ID()
  if err != nil {
    return(err)
  }
  ni.ID = foo.ID
  return ni.ipfs_addresses(foo.Addresses)
}

/*
 * local_host
 * This will check if the IP is on one of our
 * local networks 
 */
func (ni *NetInfo) local_host(ip net.IP) bool {
  for _, f := range(ni.IPS) {
    if f.Contains(ip) == true {
      return true
    }
  }
  return false
}

/* 
 * my_addresses
 * does what it sounds like it does. It 
 * will grab all local addresses and save
 * them. 
 */
func (ni *NetInfo) my_addresses() error {
  ifaces, err := net.Interfaces()
  if err != nil {
    return(err)
  }
  for _, i := range ifaces {
    addrs, err := i.Addrs()
    if err != nil {
      return(err)
    }
    for _, addr := range addrs {
      switch v := addr.(type) {
        case *net.IPNet:
          if v.IP.IsLoopback() != true {
            ni.IPS = append(ni.IPS, *v)
          }
      }
    }
  }
  return nil
}

/*
 * IPFSAddresses
 * We need to loop through our addresses and only add
 * announceable ones that are on our subnet and exclude
 * localhost 
 */
func (ni *NetInfo) ipfs_addresses(a []string) error {
  var foo IPFSAddr
  for _, f := range a {
    err := foo.Parse(f)
    if err != nil {
     return err
    }
    if ni.HaveIP(foo.IP) {
      ni.PublishAddrs = append(ni.PublishAddrs, f)
      ni.Port = append(ni.Port, foo.Port)
    }
  }
  return nil
}

/*
 * HaveIP
 * This is a utility function to see if we have an 
 * address in our network (in other words exclude
 * public NAT'd addresses
 */
func (ni *NetInfo) HaveIP(ip net.IP) bool {
  for _,i := range ni.IPS {
    if(i.IP.IsLoopback() != true) && (i.IP.Equal(ip) == true) {
      return true
    }
  }
  return false
}


/*
 * AddHost
 * This will check to see if the address is on one of
 * our local subnets. If so, it will then call IPFS
 * swarm connect so we can connect to the host
 */
func (ni *NetInfo) AddHost(h IPFSAddr) (error) {

  if (ni.local_host(h.IP) == true) && (ni.HaveIP(h.IP) == false) {
    log.Print("Adding: ", h.URL)
    sh := shell.NewShell("localhost:5001")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()
    err := sh.SwarmConnect(ctx,  h.URL)
    if err != nil {
      log.Fatal(err)
      return err
    }
  }
  return nil
}
