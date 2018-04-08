Golinks
=======

Golinks is a tiny url shortnet server that uses sqlite database. This is purely
meant for my personal use and you are more than welcome to use it too.

Build
=====

First, you need to setup all your dependencies to compile this program. I am not
most familiar with Golang tooling, I am going to just use the one provided by Golang.

```bash
$ go get -u github.com/mattn/go-sqlite3
$ go get -u github.com/jinzhu/gorm
$ go get -u github.com/gin-gonic/gin
```

After you have downloaded all the dependencies, you can now build and run it.

```bash
$ git clone https://github.com/maxking/golinks
$ cd golinks
$ go build
```

Running
=======

Just execute the binary to start the running server:

```bash
$ ./golinks
```
