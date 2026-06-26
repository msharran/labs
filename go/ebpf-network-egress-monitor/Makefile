APP=egress-monitor
GO=go
ORB=orb

.PHONY: build clean run-%

build:
	$(ORB) $(GO) build -o $(APP) .

run-%: build
	$(ORB) sudo ./$(APP) $*

clean:
	rm -f $(APP)
