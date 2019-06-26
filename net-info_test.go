package main

import (
  "testing"
  "net"
)


func TestNetInfoSetup(t *testing.T) {
  var foo NetInfo

  err := foo.Setup()
  if err != nil {
    t.Log(err)
    t.Error("Failed to setup")
  }
}


func TestNetInfoHaveIP(t *testing.T) {
  var foo = net.ParseIP("8.8.4.4")
  var bar NetInfo
  err := bar.Setup()
  if err != nil {
    t.Log(err)
    t.Error("Failed to run Setup!")
  }
  if bar.HaveIP(foo) {
    t.Error("Failed test, we don't have the IP 8.8.4.4!")
  }
  foo = net.ParseIP("127.0.0.1")
  if bar.HaveIP(foo) {
    t.Error("Failed test, we don't have the IP 8.8.4.4!")
  }
}

