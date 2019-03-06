PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

proto: $(GO)
    %.pb.go: %.proto
	    protoc --gofast_out=. $<

clean:
	rm *.pb.go
