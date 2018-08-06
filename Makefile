GOLANG_VERSION = 1.8

build: ia-mcts ia-codingame game
	@echo "Done"
battle: build
	echo '' > /tmp/wins
	BATTLE=100 ; cd bin && for i in $$(seq $$BATTLE) ; do echo "Battle $$i" ; ./game $$PWD/ia-mcts $$PWD/ia-codingame | tee -a /tmp/wins | grep 'Winner' ; done ; echo $$(grep 'Winner is Player 1' /tmp/wins | wc -l) / $$BATTLE 
test:
	docker run --rm -v $$PWD/src:/go/src/legend_of_code/ -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/legend_of_code/ && go get && go run *.go'
game:
	docker run --rm -v $$PWD/src/game.go:/go/src/game/game.go \
	-v $$PWD/bin:/go/bin  \
	-it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/game/ && go get && go build *.go'
ia-mcts:
	docker run --rm -v $$PWD/src/ia_mcts.go:/go/src/ia-mcts/ia_mcts.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-mcts/ && go get && go build *.go'

generate:
	docker run --rm -v $$PWD/src/game_generation.go:/go/src/game-generation/game_generation.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/game-generation/ && go get && go build *.go'
	./bin/game-generation

mcts:
	rm *dot *png || true
	docker run --rm -v $$PWD/src/ia_mcts.go:/go/src/ia-mcts/ia_mcts.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-mcts/ && go get && go build *.go'
	./bin/ia-mcts
	@#for f in $$(ls *dot) ; do dot -Tpng $$f -o $$f.png ; done ; rm *dot || true
ia-dummy:
	docker run --rm -v $$PWD/src/ia_dummy.go:/go/src/ia-dummy/ia_dummy.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-dummy/ && go get && go build *.go'

ia-codingame:
	docker run --rm -v $$PWD/src/ia_codingame.go:/go/src/ia-codingame/ia_codingame.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-codingame/ && go get && go build *.go'
lib:
	docker run --rm -v $$PWD/src/mcts.go:/go/src/mcts/mcts.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/mcts/ && go get && go build *.go'
main:
	docker run --rm -v $$PWD/src/main.go:/go/src/main/main.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/main/ && go get && go build *.go'
	cd bin && ./main