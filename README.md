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
    'http://localhost:8090/api/v1/cvs/add?cellId=20007114&version=1&branch=main&lockKey=key1'
```

```json
{"Code":0,"Data":{"CellId":"20007114","Version":1,"Namespace":"main","LockKey":"key1","Comment":""},"Msg":"add new version ok"}
```


# download a cell version

```shell
curl --location --request GET 'http://localhost:8090/api/v1/query/download?cellId=20007114&version=1&branch=main'

```

# lock a cell

```shell
curl --location --request POST 'http://localhost:8090/api/v1/cvs/lock' \
--data-raw '{
    "cellId" :20007114,
    "branch" :"main",
    "LockKey" : "user1",
    "lockDuration": {
      "seconds": 300
    }
}'
```

# unlock a cell

```shell
curl --location --request POST 'http://localhost:8090/api/v1/cvs/unlock' \
--data-raw '{
    "cellId" :20007114,
    "branch" :"main",
    "LockKey" : "user1"
}'
```

# find cell status

```shell

curl --location --request GET 'http://localhost:8090/api/v1/query/status?cellId=20007114&version=1&branch=main'
```

# find waiting to compile queue

```shell
  query all waiting to compile queue
  curl --location --request GET 'http://localhost:8090/api/v1/admin/compileQueueSize'

  query waiting to compile queue by branch
  curl --location --request GET 'http://localhost:8090/api/v1/admin/compileQueueSizeByBranch?branch=main'
```


# swagger

http://localhost:8090/swagger/index.html
swag init -g ./cmd/main.go -o cmd/docs
http://localhost:8090/swagger/index.html
swag init

# build

under project root dir
```shell
build/.build.sh
```
`
