GOARCH=arm GOARM=5 go build -o ravenor ./main.go
./kill.sh
./scp.sh
./run.sh
