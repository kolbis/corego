// +build integration

package cassandradb_integration_test

import (
	"testing"

	"github.com/gocql/gocql"
	logger "github.com/kolbis/corego/logger"
	cassandradb "github.com/kolbis/corego/storage/cassandradb"
)

var cassandra cassandradb.CassandraDB

func init() {
	config := cassandradb.Config{
		Hosts:        []string{"int-k8s1:32742"},
		Username:     "thelotter",
		Password:     "123",
		KeyspaceName: "golang_test",
	}
	l, _ := logger.NewLogger("test", "file logger", logger.InfoLogLevel)
	cassandra = cassandradb.NewCassandraDB(config, &l)
}

func TestSelect(t *testing.T) {
	selectQuery := `SELECT pk, cck, data
					FROM golang_test.test_table
					LIMIT 10;`
	_, err := cassandra.Select(selectQuery, gocql.Quorum)
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	insertQuery := `INSERT INTO golang_test.test_table(pk, cck, data)
					VALUES(?, ?, ?);`
	err := cassandra.ExecuteQuery(insertQuery, gocql.Quorum, `2e`, `3e`, `4e`)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	updateQuery := `UPDATE golang_test.test_table
					SET data='test'
					WHERE pk=? AND cck=?;`
	err := cassandra.ExecuteQuery(updateQuery, gocql.Quorum, `2e`, `3e`)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	deleteQuery := `DELETE FROM golang_test.test_table
					WHERE pk=? AND cck=?;`
	err := cassandra.ExecuteQuery(deleteQuery, gocql.Quorum, `2e`, `3e`)
	if err != nil {
		t.Error(err)
	}
}
