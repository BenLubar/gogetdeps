gogetdeps
=========

`go get` dependency version pinning

How to use
----------

```bash
go get -u github.com/BenLubar/gogetdeps
cd $GOPATH/src/github.com/user/project
gogetdeps
```

Can I update my dependencies?
-----------------------------

```bash
cd $GOPATH/src/github.com/user/project
gogetdeps -undo
go get -u
gogetdeps
```
