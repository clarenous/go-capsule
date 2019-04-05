PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)
API_PB_GO = api.pb.go

proto: $(GO)
%.pb.go: %.proto
	protoc -I/usr/local/include -I. -I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	--gogo_out=plugins=grpc:. \
	--grpc-gateway_out=logtostderr=true:. $<
	sed -i.bak 's/Size_/Size/' $(API_PB_GO)
	rm $(API_PB_GO).bak

clean:
	rm *.pb.go
	rm *.pb.gw.go
