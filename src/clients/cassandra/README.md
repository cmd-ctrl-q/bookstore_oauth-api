

## Cassandra 

Bad approach
``` Go - bad approach because poor performance
// file: cassandra_db.go
var (
	cluster *gocql.ClusterConfig
)

func init() {
	// connect to the cluster. use host as env var
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() (*gocql.Session, error) {
	return cluster.CreateSession()
}

// db_repository.go - for the next approach, you dont need to create more than 1 session
	session, err := cassandra.GetSession()
	if err != nil {
		errors.NewInternalServerError(err.Error())
	}
	defer session.Close()
```

Best approach
You should create only 1 session to interact with the entire db
``` Go - great performance
// you only need one session for all db connections so you should do it this way
var (
	session *gocql.Session
)

func init() {
	// connect to the cluster. use host as env var
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
}

func GetSession() *gocql.Session {
	return session
}

// file: db_repository.go
if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(...)
```