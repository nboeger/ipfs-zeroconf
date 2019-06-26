package main

import "testing"


func TestIPFSAddrParse(t *testing.T) {
  var foo IPFSAddr
  err := foo.Parse("/ip4/127.0.0.1/tcp/4001/ipfs/QmZTv9kt1r2q4uHsk2QmqFLUGBLwvmT3BskjgiNFBYVj15")
  if err != nil {
    t.Log(err)
    t.Errorf("Error couold not parse address!")
  }
  err = foo.Parse("/ip4/127.0.0.1/tcp/abc/ipfs/QmZTv9kt1r2q4uHsk2QmqFLUGBLwvmT3BskjgiNFBYVj15")
  if err == nil {
    t.Log(err)
    t.Errorf("Error couold not parse address!")
  }
  err = foo.Parse("/ip4/cookies/tcp/4001/ipfs/QmZTv9kt1r2q4uHsk2QmqFLUGBLwvmT3BskjgiNFBYVj15")
  if err == nil {
    t.Log(err)
    t.Errorf("Error couold not parse address!")
  }
  err = foo.Parse("bla bla bla/ bla")
  if err == nil {
    t.Errorf("Error should not be able to parse this string")
  }
  var non string
  err = foo.Parse(non)
  if err == nil {
    t.Errorf("Error couold not parse Bonjour message!")
  }
}

func TestIPFSAddrParseBonjour(t *testing.T) {

  var addr = []string{"ID=QmZTv9kt1r2q4uHsk2QmqFLUGBLwvmT3BskjgiNFBYVj15", "Address=/ip4/127.0.0.1/tcp/4001/ipfs/QmZTv9kt1r2q4uHsk2QmqFLUGBLwvmT3BskjgiNFBYVj15"}
  var foo IPFSAddr

  err := foo.ParseBonjour(addr)
  if err != nil {
    t.Errorf("Error couold not parse Bonjour message!")
  }
  var bla = []string{"This is a fake string"}
  err = foo.ParseBonjour(bla)
  if err == nil {
    t.Errorf("Error couold not parse Bonjour message!")
  }
  var blob = []string{"This is a fake string", "This is another"}
  err = foo.ParseBonjour(blob)
  if err == nil {
    t.Errorf("Error couold not parse Bonjour message!")
  }
  var non []string
  err = foo.ParseBonjour(non)
  if err == nil {
    t.Errorf("Error couold not parse Bonjour message!")
  }

}

