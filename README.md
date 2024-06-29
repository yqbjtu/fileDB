This project is a versioned file storage system.    
it does store the file on local disk or oss, you can access the file form this system  

# how to run


```shell
go run cmd/main.go

```


# how to build

on linux

```shell
./build.sh

```


# upload a new cell version

create /tmp/osmdb/data dir on the machine which app runs
```azure
curl -X POST -F "file=@/Users/ericyang/Downloads/20007114.osm" \
    'http://localhost:8090/api/v1/cvs/add?cellId=20007114&version=1&namespace=main&lockKey=key1'
```

```json
{"Code":0,"Data":{"CellId":"20007114","Version":1,"Namespace":"main","LockKey":"key1","Comment":""},"Msg":"add new version ok"}
```

