package relationships

import (
	"context"
	"fmt"

	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// mergeNodeFeed transaction function.
func mergeNodeFeed(f *feedsv2.Feed) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeFeedTpl, map[string]interface{}{
			"feed_id":  f.Id,
			"name":     f.Name,
			"hostname": f.Hostname,
			"url":      f.Url,
			"type":     "feed",
			"uid":      uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("feed %s record didn't create: %s", f.Hostname, result.Err().Error())
	}
}

// MergeNodeFeed upserts a feed in neo4j.
func MergeNodeFeed(ctx context.Context, neoClient *neo.Neo, f *feedsv2.Feed) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(mergeNodeFeed(f))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// mergeNodeEntity transaction function.
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

// MergeNodeEntity upserts an entity in neo4j.
func MergeNodeEntity(ctx context.Context, neoClient *neo.Neo, label string, entityType string) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(mergeNodeEntity(label, entityType))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// createNodeArticle transaction function.
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
			"hostname":     article.Hostname,
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

// CreateNodeArticle creates a new article node.
func CreateNodeArticle(ctx context.Context, neoClient *neo.Neo, article *NodeArticle) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	res, err := session.WriteTransaction(createNodeArticle(article))
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// createPublishedAt published_at transaction function.
func createPublishedAt(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(publishedAtTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createAuthorOf author_of transaction function.
func createAuthorOf(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(authorOfTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createWritesFor writes_for transaction function.
func createWritesFor(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(writesForTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createHasEntity has_entity transaction function.
func createHasEntity(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(hasEntityTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createTopic in_topic transaction function.
func createTopic(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(topicTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// MergeRel upserts a relationship between nodes.
func MergeRel(ctx context.Context, neoClient *neo.Neo, source, dest, rel string) error {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	var f neo4j.TransactionWork

	switch rel {
	case "PUBLISHED_AT":
		f = createPublishedAt(source, dest)
	case "AUTHOR_OF":
		f = createAuthorOf(source, dest)
	case "WRITES_FOR":
		f = createWritesFor(source, dest)
	case "HAS_ENTITY":
		f = createHasEntity(source, dest)
	case "IN_TOPIC":
		f = createTopic(source, dest)
	}

	if _, err := session.WriteTransaction(f); err != nil {
		return err
	}
	return nil
}

// CreateSimilar creates a neo4j relationship between source and target articles.
func CreateSimilar(ctx context.Context, neoClient *neo.Neo, source, dest string, score float64) error {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	if _, err := session.WriteTransaction(createSimilar(source, dest, score)); err != nil {
		return err
	}
	return nil
}

// createSimilar similar_with transaction function.
func createSimilar(source string, dest string, score float64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(similarTxFunc, map[string]interface{}{"source": source, "dest": dest, "score": score})
	}
}
