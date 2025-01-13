package relationships

import (
	"context"
	"fmt"

	feedsv2 "github.com/cvcio/mediawatch/pkg/mediawatch/feeds/v2"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// mergeNodeFeed transaction function.
func mergeNodeFeed(ctx context.Context, f *feedsv2.Feed) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, nodeFeedTpl, map[string]interface{}{
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

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("feed %s record didn't create: %s", f.Hostname, result.Err().Error())
	}
}

// MergeNodeFeed inserts or updates a feed in neo4j.
func MergeNodeFeed(ctx context.Context, neoClient *neo.Neo, f *feedsv2.Feed) (string, error) {
	session := neoClient.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = session.Close(ctx) }()

	res, err := session.ExecuteWrite(ctx, mergeNodeFeed(ctx, f))
	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// mergeNodeEntity transaction function.
func mergeNodeEntity(ctx context.Context, label string, entityType string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
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

		result, err := tx.Run(ctx, template, map[string]interface{}{
			"label": label,
			"type":  entityType,
			"uid":   uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("%s %s record didn't create: %s", entityType, label, result.Err().Error())
	}
}

// MergeNodeEntity inserts or updates an entity in neo4j.
func MergeNodeEntity(ctx context.Context, neoClient *neo.Neo, label string, entityType string) (string, error) {
	session := neoClient.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = session.Close(ctx) }()

	res, err := session.ExecuteWrite(ctx, mergeNodeEntity(ctx, label, entityType))
	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// createNodeArticle transaction function.
func createNodeArticle(ctx context.Context, article *NodeArticle) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, nodeArticleTpl, map[string]interface{}{
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

		// if result.Next(ctx) {
		// 	return result.Record().Values[0], nil
		// }

		// return nil, fmt.Errorf("article %s record didn't create: %s", article.DocId, result.Err().Error())
		return article.DocId, nil
	}
}

// CreateNodeArticle creates a new article node.
func CreateNodeArticle(ctx context.Context, neoClient *neo.Neo, article *NodeArticle) (string, error) {
	session := neoClient.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = session.Close(ctx) }()

	res, err := session.ExecuteWrite(ctx, createNodeArticle(ctx, article))
	if err != nil {
		return "", err
	}

	return res.(string), nil
}

// createPublishedAt published_at transaction function.
func createPublishedAt(ctx context.Context, source string, dest string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, publishedAtTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createAuthorOf author_of transaction function.
func createAuthorOf(ctx context.Context, source string, dest string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, authorOfTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createWritesFor writes_for transaction function.
func createWritesFor(ctx context.Context, source string, dest string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, writesForTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createHasEntity has_entity transaction function.
func createHasEntity(ctx context.Context, source string, dest string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, hasEntityTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// createTopic in_topic transaction function.
func createTopic(ctx context.Context, source string, dest string) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, topicTpl, map[string]interface{}{"source": source, "dest": dest})
	}
}

// MergeRel inserts or updates a relationship between nodes.
func MergeRel(ctx context.Context, neoClient *neo.Neo, source, dest, rel string) error {
	session := neoClient.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = session.Close(ctx) }()

	var f neo4j.ManagedTransactionWork

	switch rel {
	case "PUBLISHED_AT":
		f = createPublishedAt(ctx, source, dest)
	case "AUTHOR_OF":
		f = createAuthorOf(ctx, source, dest)
	case "WRITES_FOR":
		f = createWritesFor(ctx, source, dest)
	case "HAS_ENTITY":
		f = createHasEntity(ctx, source, dest)
	case "IN_TOPIC":
		f = createTopic(ctx, source, dest)
	}

	if _, err := session.ExecuteWrite(ctx, f); err != nil {
		return err
	}

	return nil
}

// CreateSimilar creates a neo4j relationship between source and target articles.
func CreateSimilar(ctx context.Context, neoClient *neo.Neo, source, dest string, score float64) error {
	session := neoClient.Client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = session.Close(ctx) }()

	if _, err := session.ExecuteWrite(ctx, createSimilar(ctx, source, dest, score)); err != nil {
		return err
	}
	return nil
}

// createSimilar similar_with transaction function.
func createSimilar(ctx context.Context, source string, dest string, score float64) neo4j.ManagedTransactionWork {
	return func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, similarTxFunc, map[string]interface{}{"source": source, "dest": dest, "score": score})
	}
}
