Golinks
=======

Golinks is a tiny url shortnet server that uses sqlite database. This is purely
meant for my personal use and you are more than welcome to use it too.

Build
=====

First, you need to setup all your dependencies to compile this program. This project
uses [dep](https://github.com/golang/dep) for dependency management.
```bash
$ go get -u github.com/golang/dep
```

Now you can clone the repo and download all the dependencies.

```bash
$ git clone https://github.com/maxking/golinks
$ cd golinks
$ dep ensure
$ go build
```

Running
=======

Just execute the binary to start the running server:

```bash
$ sudo PORT=80 ./golinks
```

Environment variable `PORT` can be used to customize the port we are listening to.


I also add an entry in `/etc/hosts` to point `go` to IP address of the machine
running Golinks, so that you can use `http://go/short` in your browser to re-direct to
the URL you want to point to.