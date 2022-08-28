build:
	cd plain/tcp-server & go build -o ../../bin/
	cd plain/tcp-client & go build -o ../../bin/
	cd tls/keygen & go build -o ../../bin/
	cd tls/tls-server & go build -o ../../bin/
	cd tls/tls-client & go build -o ../../bin/
