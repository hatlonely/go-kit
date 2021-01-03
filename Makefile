test:
	go test -gcflags=all=-l -cover ./alics
	go test -gcflags=all=-l -cover ./bind
	go test -gcflags=all=-l -cover ./cast
	go test -gcflags=all=-l -cover ./config
	go test -gcflags=all=-l -cover ./flag
	go test -gcflags=all=-l -cover ./logger
	go test -gcflags=all=-l -cover ./refx
	go test -gcflags=all=-l -cover ./rpcx
	go test -gcflags=all=-l -cover ./strx
	go test -gcflags=all=-l -cover ./validator
	go test -gcflags=all=-l -cover ./ops
