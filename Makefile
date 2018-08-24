GOLANG_VERSION = 1.8

build: gen-ia ia-codingame ia-mcts game
	@echo "Done"

test:
	docker run --rm -v $$PWD/src/game:/go/src/game \
	-v $$PWD/src/agents:/go/src/agents \
	-v $$PWD/bin:/go/bin  \
	-it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/game/ && go get && go test -v'

main: test
	docker run --rm -v $$PWD/src/game:/go/src/game \
	-v $$PWD/src/agents:/go/src/agents \
	-v $$PWD/src/ai:/go/src/ai \
	-v $$PWD/bin:/go/bin  \
	-it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ai/ && go get && go build main.go'
codingame: test
	./merge.rb

gen-ia: codingame
	docker run --rm -v $$PWD/src/gen-codingame.go:/go/src/gen-codingame/gc.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/gen-codingame/ && go get && go build *.go'


test-battle: main
	BATTLE=2; echo '' > /tmp/wins ; for i in $$(seq $$BATTLE) ; do echo -n "[$$i] " ; ./bin/ai 2>&1 | tee -a /tmp/wins | grep 'winner' ; done ; WINS=$$(grep 'player 1' /tmp/wins | wc -l) ; echo "$$WINS/$$BATTLE" 

battle: build
	echo '' > /tmp/wins
	BATTLE=1 ; cd bin && for i in $$(seq $$BATTLE) ; do echo "Battle $$i" ; ./game $$PWD/gen-codingame $$PWD/ia-codingame ; done
battle-n: build
	echo '' > /tmp/wins
	BATTLE=100 ; cd bin && for i in $$(seq $$BATTLE) ; do echo "Battle $$i" ; ./game $$PWD/gen-codingame $$PWD/ia-codingame | tee -a /tmp/wins | tail -n 2 ; done ; echo $$(grep 'Winner is Player 1' /tmp/wins | wc -l) / $$BATTLE

game:
	docker run --rm -v $$PWD/src/game.go:/go/src/game/game.go \
	-v $$PWD/bin:/go/bin  \
	-it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/game/ && go get && go build *.go'

ia-mcts:
	docker run --rm -v $$PWD/src/ia_mcts.go:/go/src/ia-mcts/ia_mcts.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-mcts/ && go get && go build *.go'
ia-dummy:
	docker run --rm -v $$PWD/src/ia_dummy.go:/go/src/ia-dummy/ia_dummy.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-dummy/ && go get && go build *.go'

ia-codingame:
	docker run --rm -v $$PWD/src/ia_codingame.go:/go/src/ia-codingame/ia_codingame.go -v $$PWD/bin:/go/bin -it golang:$(GOLANG_VERSION) bash -c 'cd /go/src/ia-codingame/ && go get && go build *.go'

png:
	@cd games ; for f in $$(ls *dot) ; do dot -Tpng $$f -o $$f.png ; done ; rm *dot || true
