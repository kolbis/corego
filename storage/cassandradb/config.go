package cassandradb

import "github.com/gocql/gocql"

// Configuration for cassandradb client.
type Config struct {

	// List of IP addresses and ports used to connect to the cluster.  For example: []string{"host1:123", "host2:124"}
	Hosts []string

	//  Username for connecting to Cassandra hosts. For example: "thelotter".
	Username string

	// Password for connecting to Cassandra hosts.
	Password string

	// A keyspace in Cassandra is a namespace that defines data replication on nodes. A cluster contains one keyspace per node.
	KeyspaceName string
}

// NewConfig creates a new Cassandra connection config
func NewConfig(hosts []string, username string, password string, keyspaceName string, port int, consistencyLevel gocql.Consistency) Config {
	return Config{
		Hosts:        hosts,
		Username:     username,
		Password:     password,
		KeyspaceName: keyspaceName,
	}
}
