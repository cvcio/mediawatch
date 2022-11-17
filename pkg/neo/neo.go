package neo

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

// Neo Neo4j Driver
type Neo struct {
	Client neo4j.Driver
}

// NewNeo Client from given config
func NewNeo(bolt, user, password string) (*Neo, error) {
	c := func(conf *neo4j.Config) {}
	driver, err := neo4j.NewDriver(bolt, neo4j.NoAuth(), c)
	if err != nil {
		return nil, err
	}

	return &Neo{Client: driver}, nil
}
