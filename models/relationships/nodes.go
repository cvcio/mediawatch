package relationships

import (
	"context"
	"fmt"
	"strings"

	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func getEntityType(entityType string) string {
	switch strings.ToLower(entityType) {
	case "feed":
		return "Feed"
	case "gpe":
		return "GPE"
	case "org":
		return "Organization"
	case "person":
		return "Person"
	case "author":
		return "Author"
	case "topic":
		return "Topic"
	default:
		return "Article"
	}
}

var nodeFeedTpl = `
	MERGE (n:Feed {
		feed_id: $feed_id,
		name: $name, 
		screen_name: $screen_name,
		url: $url,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

func mergeNodeFeed(f *feed.Feed) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeFeedTpl, map[string]interface{}{
			"feed_id":     f.ID.Hex(),
			"name":        f.Name,
			"screen_name": f.ScreenName,
			"url":         f.URL,
			"type":        "feed",
			"uid":         uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("feed %s record didn't create: %s", f.ScreenName, result.Err().Error())
	}
}

func MergeNodeFeed(ctx context.Context, neoClient *neo.Neo, f *feed.Feed) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(mergeNodeFeed(f))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

var nodeEntityTpl = `
	MERGE (n:Entity {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`
var nodeOrganizationTpl = `
	MERGE (n:Organization {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`
var nodeGPETpl = `
	MERGE (n:GPE {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`
var nodePersonTpl = `
	MERGE (n:Person {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`
var nodeTopicTpl = `
	MERGE (n:Topic {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`
var nodeAuthorTpl = `
	MERGE (n:Author {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

func mergeNodeEntity(label string, entityType string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		template := nodeEntityTpl
		if getEntityType(entityType) == "Organization" {
			template = nodeOrganizationTpl
		} else if getEntityType(entityType) == "GPE" {
			template = nodeGPETpl
		} else if getEntityType(entityType) == "Person" {
			template = nodePersonTpl
		} else if getEntityType(entityType) == "Topic" {
			template = nodeTopicTpl
		} else if getEntityType(entityType) == "Author" {
			template = nodeAuthorTpl
		}

		result, err := tx.Run(template, map[string]interface{}{
			"label": label,
			"type":  entityType,
			"uid":   uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("%s %s record didn't create: %s", entityType, label, result.Err().Error())
	}
}

func MergeNodeEntity(ctx context.Context, neoClient *neo.Neo, label string, entityType string) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(mergeNodeEntity(label, entityType))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

var nodeArticleTpl = `
	MERGE (n:Article {
		uid: $uid, 
		doc_id: $doc_id,
		lang: $lang,
		crawled_at: datetime($crawled_at),
		url: $url,
		title: $title,
		published_at: datetime($published_at),
		screen_name: $screen_name
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

func createNodeArticle(article *NodeArticle) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeArticleTpl, map[string]interface{}{
			"uid":          article.DocId,
			"doc_id":       article.DocId,
			"lang":         article.Lang,
			"crawled_at":   article.CrawledAt,
			"url":          article.Url,
			"title":        article.Title,
			"published_at": article.PublishedAt,
			"screen_name":  article.ScreenName,
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("article %s record didn't create: %s", article.DocId, result.Err().Error())
	}
}

func CreateNodeArticle(ctx context.Context, neoClient *neo.Neo, article *NodeArticle) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(createNodeArticle(article))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

var publishedAtTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:PUBLISHED_AT]->(b)
`

func CreatePublishedAtTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(publishedAtTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

var authorOfTpl = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Article {uid: $dest})
	MERGE (a)-[:AUTHOR_OF]->(b)
`

func CreateAuthorOfTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(authorOfTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

var writesForTpl = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:WRITES_FOR]->(b)
`

func CreateWritesForTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(writesForTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

var hasEntityTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Entity {uid: $dest})
	MERGE (a)-[:HAS_ENTITY]-(b)
`

func CreateHasEntityTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(hasEntityTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

var topicTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Topic {uid: $dest})
	MERGE (a)-[:IN_TOPIC]-(b)
`

func CreateTopicTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(topicTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}
func MergeRel(ctx context.Context, neoClient *neo.Neo, source, dest, rel string) error {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	var f neo4j.TransactionWork

	switch rel {
	case "PUBLISHED_AT":
		f = CreatePublishedAtTxFunc(source, dest)
	case "AUTHOR_OF":
		f = CreateAuthorOfTxFunc(source, dest)
	case "WRITES_FOR":
		f = CreateWritesForTxFunc(source, dest)
	case "HAS_ENTITY":
		f = CreateHasEntityTxFunc(source, dest)
	case "IN_TOPIC":
		f = CreateTopicTxFunc(source, dest)
	}

	if _, err := session.WriteTransaction(f); err != nil {
		return err
	}
	return nil
}
