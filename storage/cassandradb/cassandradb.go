package cassandradb

import (
	"github.com/gocql/gocql"
	tleerrors "github.com/kolbis/corego/errors"
	tlelogger "github.com/kolbis/corego/logger"
)

// CassandraDB interface defines methods to interact with Cassandra database
type CassandraDB interface {
	Select(cql string, consistencyLevel gocql.Consistency, params ...interface{}) ([]map[string]interface{}, error)
	ExecuteQuery(cql string, consistencyLevel gocql.Consistency, params ...interface{}) error
	CloseSession()
}

// cassandradb struct it's used to implement CassandraDB interface
type cassandradb struct {
	Session *gocql.Session
	Logger  *tlelogger.Logger
	Error   error
}

// KeyspaceSessionDictionary contains pairs of keyspaceName - cassandra session
var KeyspaceSessionDictionary = map[string]*gocql.Session{}

// NewCassandraDB function creates a session to CassandraDB and exposes methods to interact with it
// A session is similar to a connection pool in sql.
// When the application starts, create a session and keep a reference to it.
// When the application shutdown, make sure you close the session
func NewCassandraDB(config Config, logger *tlelogger.Logger) CassandraDB {
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.Username,
		Password: config.Password,
	}
	cluster.Keyspace = config.KeyspaceName
	session, err := getOrCreateSession(cluster, config.KeyspaceName, *logger)

	return &cassandradb{
		Session: session,
		Logger:  logger,
		Error:   err,
	}
}

// Select returns the data from the query in the form of []map[string]interface{}
// consistencyLevel = Level of consistency (ex: gocql.Quorum)
// cql = CQL string, for example: "WHERE id = ?"
// params = Values of positional args
func (db *cassandradb) Select(cql string, consistencyLevel gocql.Consistency, params ...interface{}) ([]map[string]interface{}, error) {
	if db.Session == nil {
		errDatabase := tleerrors.NewDatabaseErrorf("Cassandra session must be initialized before calling select: %s", cql)
		return nil, errDatabase
	}

	iter := db.Session.Query(cql, params...).Consistency(consistencyLevel).Iter()
	result, err := iter.SliceMap()
	defer iter.Close()

	if err != nil {
		logDatabaseError(*db.Logger, err, "Error on executing iter.SliceMap() for following query: "+cql)
		return nil, err
	}

	return result, nil
}

// ExecuteQuery can be used for executing insert, update and delete.
// consistencyLevel = Level of consistency (ex: gocql.Quorum)
// cql = CQL string, for example: "WHERE id = ?"
// params = Values of positional args
func (db *cassandradb) ExecuteQuery(cql string, consistencyLevel gocql.Consistency, params ...interface{}) error {
	if db.Session == nil {
		errDatabase := tleerrors.NewDatabaseErrorf("Cassandra session must be initialized before calling select: %s", cql)
		return errDatabase
	}

	err := db.Session.Query(cql, params...).Consistency(consistencyLevel).Exec()
	if err != nil {
		logDatabaseError(*db.Logger, err, "Error on executing following query: "+cql)
	}

	return err
}

// CloseSession closes CassandraDB connection.
// The session is unusable after this operation.
func (db *cassandradb) CloseSession() {
	db.Session.Close()
}

// getOrCreateSession creates a new CassandraDB session if it doesn't exists or retrieve the existing one by keyspace name
func getOrCreateSession(cluster *gocql.ClusterConfig, keyspaceName string, logger tlelogger.Logger) (*gocql.Session, error) {
	session := KeyspaceSessionDictionary[keyspaceName]
	if session != nil {
		return session, nil
	}
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		logDatabaseError(logger, err, "An error occured while trying to create a cassandra session for keyspace: "+cluster.Keyspace)
		return nil, err
	}
	KeyspaceSessionDictionary[keyspaceName] = session
	return session, err
}

// logDatabaseError logs the database error
func logDatabaseError(logger tlelogger.Logger, err error, msg string) {
	errDatabase := tleerrors.NewDatabaseError(err, msg)
	wrappedErrDatabase := tleerrors.Wrap(err, errDatabase)
	logger.Log(wrappedErrDatabase)
}
