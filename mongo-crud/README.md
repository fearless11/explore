docker install mongodb

[TOC]

[官网](https://docs.mongodb.com/manual/reference/method/#user-management-methods)
[mongodoc](https://godoc.org/go.mongodb.org/mongo-driver/mongo#example-Client-Connect)
[github](https://github.com/mongodb/mongo-go-driver)

#### 安装

```shell
docker pull mongo       
docker image ls 
docker run -p 27017:27017 -v /data/mongo:/db/mongo --name docker-mongo -d mongo  // 启动容器
docker container ls              // 查看容器
docker stop docker-mongo         // 停止容器
docker start docker-mongo        // 启动容器
docker exec -it docker-mongo mongo admin  // 进入mongo
```

#### 操作

```
help                                 
use alert                            //创建数据库
db.project.insertOne( { x: 1 } )     //创建集合collection 并插入数据
show dbs                             //查看db
show collections                     //查看集合
db.project.find()                    //查询数据
db.project.find({ x:1 }).pretty()    //指定条件查询，json显示
```

