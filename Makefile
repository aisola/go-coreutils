
all:
	go build cat.go
	go build echo.go
	go build exit.go
	go build false.go
	go build logname.go
	go build mv.go
	go build pwd.go
	go build rm.go
	go build sleep.go
	go build uptime.go
	go build whoami.go

fmt:
	go fmt cat.go
	go fmt echo.go
	go fmt exit.go
	go fmt false.go
	go fmt logname.go
	go fmt mv.go
	go fmt pwd.go
	go fmt rm.go
	go fmt sleep.go
	go fmt uptime.go
	go fmt whoami.go

clean:
	rm -f cat echo exit false logname mv pwd rm sleep uptime whoami
	rm -f *.exe
	rm -f .deps
