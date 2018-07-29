build: ia game
	@echo "Done"
test:
	docker run --rm -v $$PWD/src:/go/src/legend_of_code/ -it golang:1.8 bash -c 'cd /go/src/legend_of_code/ && go get && go run *.go'
game:
	docker run --rm -v $$PWD/src/game.go:/go/src/game/game.go \
	-v $$PWD/bin:/go/bin  \
	-it golang:1.8 bash -c 'cd /go/src/game/ && go get && go build *.go'
	cd bin && ./game $$PWD/ia-dummy $$PWD/ia-dummy
ia:
	docker run --rm -v $$PWD/src/ia_dummy.go:/go/src/ia-dummy/ia_dummy.go -v $$PWD/bin:/go/bin -it golang:1.8 bash -c 'cd /go/src/ia-dummy/ && go get && go build *.go'

