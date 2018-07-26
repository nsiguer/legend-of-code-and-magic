test:
	docker run --rm -v $$PWD/src:/go/src/legend_of_code/ -it golang:1.8 bash -c 'cd /go/src/legend_of_code/ && go get && go run *.go'

