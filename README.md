# gpm
Go Proxy Manager

# Build
```
go build
```

# usage
```
Usage of gpm:
    -address string
        listening address (default "0.0.0.0:3000")
    -file string
        path to config file (default "gpm.toml")
    -host value
        from->to, ex: gh.localhost:3000--https://github.com
    -static value
        repo^branch^name, ex: github.com/simba-fs/gpm^main^blog
    -storage string
        directory to store files such as static files (default "./storage")
```

# TODO 
- [ ] Check if the file exists when setting up a static host
