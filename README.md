# DnsCompare
![workflow](https://github.com/DBHeise/DnsCompare/actions/workflows/main.yml/badge.svg)
[![Build Status](https://travis-ci.org/DBHeise/DnsCompare.svg?branch=master)](https://travis-ci.org/DBHeise/DnsCompare)
[![Go Report Card](https://goreportcard.com/badge/github.com/DBHeise/DnsCompare)](https://goreportcard.com/report/github.com/DBHeise/DnsCompare)
[![Inline docs](http://inch-ci.org/github/DBHeise/DnsCompare.svg?branch=master)](http://inch-ci.org/github/DBHeise/DnsCompare)

A tool for comparing DNS results from multiple DNS servers.

## Why?

Could this be done with a bash script and dig, or powershell and nslookup? yes (we'll almost). Why do you want to do this though? It has a couple of purposes:

1. You can compare your dns record propogation
2. You can check for domains that are being blocked/filtered and/or using GeoDNS (or something simliar)

## How to install

## How to use

DnsCompare has several commands that it can use. 

### Compare

Simple Example
```shell
DnsCompare.exe compare -t example.com

AdGuard (Default)                [93.184.216.34]
AdGuard (Family Protection)      [93.184.216.34]
AdGuard (Non-Filtering)          [93.184.216.34]
CleanBrowsing (Adult Filter)     [93.184.216.34]
CleanBrowsing (Family Filter)    [93.184.216.34]
CleanBrowsing (Security Filter)  [93.184.216.34]
Cloudflare                       [93.184.216.34]
Comodo (Secure DNS)              [93.184.216.34]
Comodo (Secure Internet Gateway) [93.184.216.34]
CyberGhost                       [93.184.216.34]
DNS.Watch                        [93.184.216.34]
DnsForge.de                      [93.184.216.34]
FreeDNS                          [93.184.216.34]
Freenom World                    [93.184.216.34]
Freifunk München e.V.            [93.184.216.34]
French Data Network (FDN)        [93.184.216.34]
Google                           [93.184.216.34]
Hurricane Electric               [93.184.216.34]
Lightning Wire Labs              [93.184.216.34]
Neustar DNS Advantage            [93.184.216.34]
NextDNS                          [93.184.216.34]
OpenDNS                          [93.184.216.34]
OpenDNS (FamilyShield)           [93.184.216.34]
Quad9 (Secured w/ECS)            [93.184.216.34]
Quad9 (Secured)                  [93.184.216.34]
Quad9 (Unsecured)                [93.184.216.34]
Securolytics                     [93.184.216.34]
Sprintlink General DNS           [93.184.216.34]
UncensoredDNS                    [93.184.216.34]
Verisign                         [93.184.216.34]
Yandex (Basic)                   [93.184.216.34]
Yandex (Family)                  [93.184.216.34]
Yandex (Safe)                    [93.184.216.34]
dns.sb                           [93.184.216.34]
puntCat                          [93.184.216.34]
```

Complex (command line) Example
```shell
DnsCompare.exe compare -m dot -p -k random4 -r A -t example.com

CMRG                                                [93.184.216.34]
CleanBrowsing (Adult Filter)                        [93.184.216.34]
CleanBrowsing (Family Filter)                       [93.184.216.34]
CleanBrowsing (Security Filter)                     [93.184.216.34]
Digital Courage                                     [93.184.216.34]
Digitale Gesellschaft (Digital Soceity Switzerland) [93.184.216.34]
Dismail.de                                          [93.184.216.34]
DnsForge.de                                         [93.184.216.34]
Foundation for Applied Privacy                      [93.184.216.34]
Freifunk München e.V.                               [93.184.216.34]
GetDNS                                              [93.184.216.34]
Google                                              [93.184.216.34]
Post-Factum                                         [93.184.216.34]
Quad9 (Secured w/ECS)                               [93.184.216.34]
Quad9 (Secured)                                     [93.184.216.34]
Quad9 (Unsecured)                                   [93.184.216.34]
Restena Foundation                                  [93.184.216.34]
Snopyta                                             [93.184.216.34]
UncensoredDNS                                       [93.184.216.34]
dns.sb                                              [93.184.216.34]

```

Avaible Options:
- DNSMode - which mode of DNS query to use
  - "--DNSMode" or "-m"
  - Available strings
    - udp - tradtional UDP dns lookup
    - tcp - DNS over TCP
    - dot - DNS over TLS
    - doh - DNS over HTTP(s)
- Parallel - run queries in parallel
  - "--Parallel" or "-p"
- Pick Method - how to choose a DNS server when multiple ones are available
  - "--PickMethod" or "-k"
  - Available strings
    - first - pick the first one in the list
    - first4 - pick the first one in the list that is an IPv4 address
    - first6 - pick the first one in the lsit that is an IPv6 address
    - random - pick one at random
    - random4 - pick one at random that is an IPv4 address
    - random6 - pick one at random that is an IPv6 address
    - all - use all the available addresses
    - all4 - use all the available IPv4 addresses
    - all6 - us all the available IPv6 addresses
- RecordType - what type of DNS record are we asking for
  - "--RecordType" or "-r"
  - Available Options (more options could be added, as wanted/needed/requested)
    - A
    - AAAA
    - NS
    - CNAME
    - PTR
    - TXT
    - SRV
    - SOA
    - SIG
- Server Defition File - the json file that defines all the servers to use
  - "--ServerDef" or "-s"
  - JSON file that contains all the definitions: see ()
- Target - the DNS query
  - "--target" or "-t"
  - the domain to resolve


### Verify
verifies all the entries in the server definition file

### Version
shows the current build/version

### Common Global Command-Line Options

```
--Log.Colors              Use colors in logging (only valid for console logging) (default true)
--Log.Format string       Log format (available options: plain, json) (default "info")
--Log.Level string        Log Level (available options: panic,fatal,error,warn,info,debug,trace) (default "info")
--Log.Output string       Log output (available options: console, {filename}) (default "console")
--Log.ReportCaller        Report the calling function to the log
--Log.TimeFormat string   Log timestamp format (available options: ansi=ansic,unix=unixdate,ruby=rubydate,rfc822,rfc822z,rfc850,rfc1123,rfc1123z,json=rfc3339,rfc3339nano,kitchen,stamp,stampmili,stampmicro,stampnano,{golang time format string}) (default "RFC3339")
```


### Bugs & Requests

bug reports, and feature requests are welcome

### Contributors
![Contributors](https://contrib.rocks/image?repo=DBHeise/DnsCompare)
