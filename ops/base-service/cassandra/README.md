# cassandra

## 启动`cassandra`环境

```sh
docker run --name cassandra -p 7000:7000 -p 7001:7001 -p 7199:7199 -p 9042:9042 -p 9160:9160 -v ${PWD}/data:/var/lib/cassandra -d cassandra:3.11
```

## 创建`keyspace`

```sh
docker exec -it cassandra cassandra-cli
create keyspace example
```

##  创建tweet数据表scheme

- 运行cqlsh 客户端

```sh
curl -o cassandra-3.11.4-bin.tar.gz http://mirror.bit.edu.cn/apache/cassandra/3.11.4/apache-cassandra-3.11.4-bin.tar.gz
tar zxvf apache-cassandra-3.11.4-bin.tar.gz
cd apache-cassandra-3.11.4-bin && ./bin/cqlsh
```

```sh

CREATE TABLE tweet (
    id int PRIMARY KEY,
    timeline text,
    text text
);
```

## 写数据

```go
package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a tweet
	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", 1, "hello world").Exec(); err != nil {
		log.Fatal(err)
	}

}
```

## 读数据

```go
package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	var id int
	var text string

	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE id = ?`, 1).Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err, 123)
	}
}

```