# muescheli poc

small service that connects to a ClamAV daemon through tcp

## test locally

build ClamAV image and run
```bash
$ docker build -t clamav docker
$ docker run -p 3310:3310 clamav
```

run the webservice (requires a go installation with dep)
```bash
$ dep ensure
$ go run main.go
```

run a test
```bash
$ curl -v -X POST -H "Content-Type: multipart/form-data" -F "file1=@test/eicar.com" -F "file2=@test/test.txt"  http://localhost:8091/scan
```
or
```bash
$ curl -v -X PUT --data-binary @test/eicar.com http://localhost:8091/scan
```
or
```bash
$ curl -v http://localhost:8091/scan?url=https://tech.ebu.ch/docs/tech/tech3285.pdf
```

## run on kubernetes

if you have the monostream [helm chart](https://github.com/monostream/helm-charts) repo installed
```bash
$ helm install monostream/muescheli
```

or manually deploy
```bash
$ kubectl create -f k8s/deployment.yml
```
this will start a pod with 2 containers (clamav and muescheli)

muescheli service exposed through random NodePort port