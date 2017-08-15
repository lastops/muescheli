# muescheli poc

small service that connects to a ClamAV daemon through tcp

## test locally

build ClamAV image and run
```bash
$ docker build -t clamav docker
$ docker run -p 3310:3310 clamav
```

run the webservice (requires a go installation)
```bash
$ go run main.go
```

run a test
```bash
$ curl -i -X POST -H "Content-Type: multipart/form-data" -F "file1=@test/eicar.com" -F "file2=@test/test.txt"  http://localhost:8091/scan
```