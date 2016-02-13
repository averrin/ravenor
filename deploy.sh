GOARCH=arm GOARM=5 go build -o ravenor ./main.go
./scp.sh
./run.sh
