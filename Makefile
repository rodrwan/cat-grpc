# Makefile for go program
.PHONY:	all proto osx linux docker clean usage

PROGRAM="Category Micro Service"
SERVICE=categoryapi
LINUXBIN=catty
OSXBIN=catty-osx
IMG=catty

all: usage

proto p:
	protoc -I ./codec ./codec/$(SERVICE).proto --go_out=plugins=grpc:$(GOPATH)/src/
	@ls -al ./codec

osx o:
	@rm $(OSXBIN) | go build -a -o $(OSXBIN) server.go
	go build -o client/client client/client.go

linux l:
	@rm $(LINUXBIN) | GOOS=linux go build -a -o $(LINUXBIN) \
	--ldflags '-extldflags "-static"' \
	-tags netgo -installsuffix netgo \
	./server.go

docker d:
	docker build -t $(IMG) .

clean:
	@rm -rf ./categoryapi
	@rm client/client
	@rm $(LINUXBIN)
	@rm $(OSXBIN)

usage:
	@echo ""
	@echo "Makefile for $(PROGRAM)"
	@echo ""
	@echo "make [proto|osx|linux|docker|clean]"
	@echo "   - proto  : compile interface spec"
	@echo "   - osx    : compile client/server"
	@echo "   - linux  : compile server for linux arch"
	@echo "   - docker : create docker image"
	@echo "   - clean  : clean files"
	@echo ""
