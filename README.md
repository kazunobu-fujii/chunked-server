# chunked-server

Transfer-Encoding: chunked server

# Usage

```sh
$ go run main.go -t "application/json; charset=utf-8" -d response.dat
```

# Command Line Options

### -c [size]

chunk size (default 8)

### -d [filename]

response body (default "response.dat")

### -disable

disable chunk mode

### -t [content-type]

content type (default "text/html; charset=UTF-8")

### -s [server]

listening server (default "localhost:8080")

### -w [delay]

chunk delay (ms) (default 10)
