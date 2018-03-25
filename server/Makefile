all: build certificates

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

certificates:
	cp /etc/ssl/certs/ca-certificates.crt .

clean:
	rm ./main
	rm ./ca-certificates.crt
