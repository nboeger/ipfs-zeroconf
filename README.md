# ipfs-zeroconf

This is a Bonjour/Zeroconf utility for IPFS to add peers in the local network. 


## Install

This requires golang v1.13 or higher! 
 
Clone the repo or you can run:
```shell
  
  $ go get -u github.com/nboeger/ipfs-zeroconf 
 
```

## Running 

To run the utility just build and run `./ipfs-zeroconf`

## Options

Below are the command line options

- `-d` Turn on debugging
- `-t <int>` Timeout
- `-s <string>` Service string to be used when announcing
- `-D <string>` Domain to announced on (default is local)
- `-cmd <command>` Right now its only stop to shutdown

