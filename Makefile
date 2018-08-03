build: ia-mcts ia-codingame game
	@echo "Done"
test:
	docker run --rm -v $$PWD/src:/go/src/legend_of_code/ -it golang:1.8 bash -c 'cd /go/src/legend_of_code/ && go get && go run *.go'
game:
	docker run --rm -v $$PWD/src/game.go:/go/src/game/game.go \
	-v $$PWD/bin:/go/bin  \
	-it golang:1.8 bash -c 'cd /go/src/game/ && go get && go build *.go'
	cd bin && ./game $$PWD/ia-mcts $$PWD/ia-codingame
ia-mcts:
	docker run --rm -v $$PWD/src/ia_mcts.go:/go/src/ia-mcts/ia_mcts.go -v $$PWD/bin:/go/bin -it golang:1.8 bash -c 'cd /go/src/ia-mcts/ && go get && go build *.go'

mcts:
	docker run --rm -v $$PWD/src/ia_mcts.go:/go/src/ia-mcts/ia_mcts.go -v $$PWD/bin:/go/bin -it golang:1.8 bash -c 'cd /go/src/ia-mcts/ && go get && go build *.go'
	cd bin && ./ia-mcts 
ia-dummy:
	docker run --rm -v $$PWD/src/ia_dummy.go:/go/src/ia-dummy/ia_dummy.go -v $$PWD/bin:/go/bin -it golang:1.8 bash -c 'cd /go/src/ia-dummy/ && go get && go build *.go'

ia-codingame:
	docker run --rm -v $$PWD/src/ia_codingame.go:/go/src/ia-codingame/ia_codingame.go -v $$PWD/bin:/go/bin -it golang:1.8 bash -c 'cd /go/src/ia-codingame/ && go get && go build *.go'

