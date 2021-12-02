# Tool Box via Web written in Go

The goal of this project is to generate a **SINGLE STATIC BINARY** to serve 
as a collection of web tools.

## project layout

* /build/, contains all build related scripts, config files.
* /release/, contains generated binaries, excluded from git repo.
* /apps/, each app has its own folder under this.
* /router/, the web framework. All apps register theirselves to the router.
* main.go, the app main.

## How to build?

### Build a single binary
```
./build/build-app.sh
```

### Build a container image
```
./build/build-docker.sh
```
