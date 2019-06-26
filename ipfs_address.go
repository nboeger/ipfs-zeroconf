
package main


import (
  "net"
  "strings"
  "strconv"
  "errors"
)

/*
 * IPFSAddr 
 * Struct to keep our data for our address
 */
type IPFSAddr struct {
  URL         string
  Transport   string
  IP          net.IP
  Protocol    string
  Port        int
  Type        string
  ID          string
}
const ADDRLEN = 7

const IPFS_TRANSPORT = 1
const IPFS_IP_ADDR   = 2
const IPFS_PROTOCOL  = 3
const IPFS_PORT      = 4
const IPFS_TYPE      = 5
const IPFS_ID        = 6


/*
 * Parse
 * Simple utility to parse an IPFS address
 */
func (i *IPFSAddr) Parse(s string) (error) {
  var err error
  bar := strings.Split(s, "/")
  if len(bar) == ADDRLEN {
    i.URL = s
    i.Transport = bar[IPFS_TRANSPORT]
    i.Protocol  = bar[IPFS_PROTOCOL]
    i.Type      = bar[IPFS_TYPE]
    i.ID        = bar[IPFS_ID]
    i.Port,err  = strconv.Atoi(bar[IPFS_PORT])
    if err != nil {
      return  err
    }
    i.IP        = net.ParseIP(bar[IPFS_IP_ADDR])
    if i.IP == nil {
      return  errors.New("Error cannot parse IPFS IP address: wrong type!")
    }
  } else {
    return  errors.New("Error cannot parse IPFS address: wrong length!")
  }
  return nil
}

/*
 * ParseBonjour
 * This will parse the full bonjour message
 */
func (i *IPFSAddr) ParseBonjour(s []string) (error) {
  if len(s) == 2 {
    bar := strings.Split(s[1], "=")
    if len(bar) == 2 {
      return i.Parse(bar[1])
    }
  }
  return errors.New("ParseBonjour: Not a proper bonjour string!")
}
