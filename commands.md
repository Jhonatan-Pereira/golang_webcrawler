# commands
```sh
go mod init web_crawler

go mod tidy

go get -u golang.org/x/net/html

go run main.go

docker run -d --name mongo-go-webcrawler -p 27017:27017 mongo

docker exec -it mongo-go-webcrawler /bin/bash

mongo

show dbs

use crawler

show collections;

db.links.count({});
db.links.find({});
```