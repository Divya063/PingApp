## PingApp

A Ping CLI application for MacOS. The CLI app accepts a hostname or an IP address as its argument and supports both <b>IPv4</b> and <b> IPv6</b>, then send ICMP "echo requests" in a loop to the target while receiving "echo reply" messages. It reports loss and RTT times for each sent message.

### Setup

```bash
go install .

```
### Usage

```
sudo pingApp ping host [--count] [--interval] 
```
Examples:

```
# ping google continuously
$ sudo pingApp ping google.com
# ping google 5 times
$ sudo pingApp ping --count 5 google.com
# ping google 5 times at 2 seconds intervals
$ sudo pingApp ping google.com  --count 3 --interval 2
```



To run the tests use the command:
`sudo go test `

### Demo
[![asciicast](https://asciinema.org/a/6VT6StWYtUisZ7KbrLr7QMI1i.svg)](https://asciinema.org/a/6VT6StWYtUisZ7KbrLr7QMI1i)
