# TCP/UDP/IP stack

## The Goal

The end result will hopefully be a _minimal_ but functional and RFC-compliant TCP/UDP/IP stack that works from (at least for Linux) the network card to an HTTP server I can curl to

## RFCs

[IP](https://datatracker.ietf.org/doc/html/rfc791)
[ICMP](https://datatracker.ietf.org/doc/html/rfc792)
[UDP](https://datatracker.ietf.org/doc/html/rfc768)
[TCP](https://datatracker.ietf.org/doc/html/rfc9293)
[FTP](https://datatracker.ietf.org/doc/html/rfc959)
[Telnet](https://datatracker.ietf.org/doc/html/rfc854)

Will also eventually include...
- ARP
- Reading data off network card (only on Linux, probably)
- HTTP
- DNS
- DHCP

## Rationale

This is a project for me to relearn and deeply understand how TCP/UDP/IP stack works. It has no real reason to exist other than that; a learning project.

Most of this will be written in Golang and C
- Starting with Golang since I know it decently well, and want to get started
- Will eventually start writing parts in C to get a deeper understanding of low level networking concepts, as well as learn C

## Project Structure

### Core code

This will be the core, abstract code that implements the protocols.

### Commands

This will be the command line tools that use the core code to provide functionality like ping, traceroute, ftp client/server, telnet client/server, etc.

### Tooling

For now, starting with just vanilla Golang with testing and linting.

I will _eventually_ add CI/CD and Nx. Possibly also an artifact repository for funsies for reusable packages I create within this project.

#### Golang Libs

Golangci-lint
CompileDaemon

https://pkg.go.dev/golang.org/x/sys/unix
https://pkg.go.dev/golang.org/x/net
- Maybe

I will add other libraries as needed, primarily (or maybe only) for testing or generic golang utils

#### AI Usage

AI will only be used for code review, and debugging if needed.

## Concerns

Anything that touches the actual network card will only be done for Linux, at least initially. MacOS + BPF would be a great thing to learn about - _later_.

I will not be using the "normal" port numbers, to avoid needing sudo/permissions.

Note to self: MacOS firewall might activate on anything that is _not local_.

## Order

Starting with UDP since it is so simple and I need to get used to reading RFCs again. Also need to setup linting, testing and building

Then probably ICMP (ping)
