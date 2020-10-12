# go-fly
## web
> A web frame based on gin
### build example project
```bash
cd examples/web
go build -o srv main.go
./srv
```
### test cases
```bash
curl -H 'Token:xyz' 'localhost:8080/hosts?id=0&id=1'
# resp: [{"id":0,"name":"localhost"},{"id":1,"name":"localhost"}]
curl -H 'Token:xyz' 'localhost:8080/hosts/1/disks'
# resp: [{"id":2,"host_id":1,"dev":"/dev/sda"},{"id":2,"host_id":1,"dev":"/dev/sdb"}]
```