package nodes

import (
	"context"
	"fmt"
	"log"

	"github.com/cvcio/mediawatch/models/deprecated/article"
	"github.com/cvcio/mediawatch/models/deprecated/feed"
	"github.com/cvcio/mediawatch/pkg/neo"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// EnsureIndex fix the indexes in the Neo4J Database
func EnsureIndex(ctx context.Context) {
	indexCypher := `
		CREATE CONSTRAINT article_uid ON (a:Article) ASSERT a.uid IS UNIQUE;
		CREATE CONSTRAINT article_doc_id ON (a:Article) ASSERT a.doc_id IS UNIQUE;

		// CREATE CONSTRAINT article_tweetId ON (a:Article) ASSERT a.tweet_id IS UNIQUE;
		// CREATE CONSTRAINT article_tweetIdStr ON (a:Article) ASSERT a.tweet_id_str IS UNIQUE;

		CREATE CONSTRAINT feed_uid ON (f:Feed) ASSERT f.uid IS UNIQUE;
		CREATE CONSTRAINT feed_screen_name ON (f:Feed) ASSERT f.screen_name IS UNIQUE;

		// CREATE CONSTRAINT feed_twitter_id_str ON (f:Feed) ASSERT f.twitter_id_str IS UNIQUE;
		// CREATE CONSTRAINT feed_twitter_id ON (f:Feed) ASSERT f.twitter_id IS UNIQUE;

		CREATE CONSTRAINT author_uid ON (a:Author) ASSERT a.uid IS UNIQUE;
		CREATE CONSTRAINT author_name ON (a:Author) ASSERT a.author IS UNIQUE;

		CREATE CONSTRAINT entity_uid ON (e:Entity) ASSERT e.uid IS UNIQUE;
		CREATE CONSTRAINT entity_name ON (e:Entity) ASSERT e.entity_text IS UNIQUE;

		CREATE INDEX feed_meta FOR (f:Feed) ON (f.screen_name, f.lang);
		CREATE INDEX article_meta FOR (a:Article) ON (a.screen_name, a.url, a.lang);

		DROP CONSTRAINT article_tweetId IF EXISTS
		DROP CONSTRAINT article_tweetIdStr IF EXISTS
		DROP CONSTRAINT feed_twitter_id_str IF EXISTS
		DROP CONSTRAINT feed_twitter_id IF EXISTS

		CALL db.index.fulltext.createNodeIndex("text", ["Article"], ["title", "body", "keywords"]);
	`

	log.Println(indexCypher)
}

func ArticleNodeExtist(ctx context.Context, neoClient *neo.Neo, tweet_id_str string) bool {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	res, err := session.ReadTransaction(ExistsTxFunc(tweet_id_str))
	if err != nil {
		return false
	}
	return res != nil
}

var existsTxFunc = `
	MATCH (n:Article {tweet_id_str: $tweet_id_str}) 
	RETURN n IS NOT NULL AS exists
`

// ExistsTxFunc check if article exists transaction function
func ExistsTxFunc(tweet_id_str string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(existsTxFunc, map[string]interface{}{"tweet_id_str": tweet_id_str})
		if err != nil {
			return nil, err
		}
		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, nil
	}
}

// SaveRelations to dgraph if exists
func CreateSimilar(ctx context.Context, neoClient *neo.Neo, source, dest article.Document, score float64) error {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	if _, err := session.WriteTransaction(CreateSimilarTxFunc(source.DocID, dest.DocID, score)); err != nil {
		return err
	}

	return nil
}

var similarTxFunc = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Article {uid: $dest})
	MERGE (a)-[:SIMILAR_WITH { score: $score }]->(b)
`

// CreateSimilarTxFunc Create article's similarity transaction function
func CreateSimilarTxFunc(source string, dest string, score float64) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(similarTxFunc, map[string]interface{}{"source": source, "dest": dest, "score": score})
	}
}

func MergeNodeFeed(ctx context.Context, neoClient *neo.Neo, f *feed.Feed) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(MergeNodeFeedTxFunc(f))

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

var nodeFeedTxFunc = `
	MERGE (n:Feed {
		name: $name, 
		screen_name: $screen_name,
		twitter_id: $twitter_id, 
		twitter_id_str: $twitter_id_str, 
		twitter_profile_image: $twitter_profile_image,
		url: $url
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

// MergeNodeFeedTxFunc Create new Node Feed
func MergeNodeFeedTxFunc(f *feed.Feed) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeFeedTxFunc, map[string]interface{}{
			"name":                  f.Name,
			"screen_name":           f.ScreenName,
			"twitter_id":            f.TwitterID,
			"twitter_id_str":        f.TwitterIDStr,
			"twitter_profile_image": f.TwitterProfileImage,
			"url":                   f.URL,
			"uid":                   uuid.New().String(),
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

var nodeEntityTxFunc = `
	MERGE (n:Entity {
		entity_text: $entity_text, 
		entity_type: $entity_type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

// MergeNodeEntityTxFunc Create new Node Feed
func MergeNodeEntityTxFunc(entityText, entityType string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeEntityTxFunc, map[string]interface{}{
			"entity_text": entityText,
			"entity_type": entityType,
			"uid":         uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("entity %s record didn't create: %s", entityText, result.Err().Error())
	}
}

func MergeNodeAuthor(ctx context.Context, neoClient *neo.Neo, author string) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(MergeNodeAuthorTxFunc(author))

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

var nodeAuthorTxFunc = `
	MERGE (n:Author {
		author: $author
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

// MergeNodeAuthorTxFunc Create new Node Feed
func MergeNodeAuthorTxFunc(author string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(nodeAuthorTxFunc, map[string]interface{}{
			"author": author,
			"uid":    uuid.New().String(),
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("author %s record didn't create: %s", author, result.Err().Error())
	}
}

func CreateNodeArticle(ctx context.Context, neoClient *neo.Neo, article *NodeArticle) (string, error) {
	session := neoClient.Client.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	res, err := session.WriteTransaction(CreateNodeArticleTxFunc(article))

	if err != nil {
		return "", err
	}

	return res.(string), nil
}

var nodeArticleTxFunc = `
	MERGE (n:Article {
		uid: $uid, 
		docId: $docId,
		lang: $lang,
		crawledAt: datetime($crawledAt),
		url: $url,
		tweet_id: $tweet_id,
		tweet_id_str: $tweet_id_str,
		title: $title,
		summary: $summary,
		body: $body,
		tags: $tags,
		categories: $categories,
		publishedAt: datetime($publishedAt),
		editedAt: datetime($editedAt),
		keywords: $keywords,
		topics: $topics,
		screen_name: $screen_name
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid
`

// CreateNodeArticleTxFunc Create new Node Feed
func CreateNodeArticleTxFunc(article *NodeArticle) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {

		result, err := tx.Run(nodeArticleTxFunc, map[string]interface{}{
			"uid":          article.DocID,
			"docId":        article.DocID,
			"lang":         article.Lang,
			"crawledAt":    article.CrawledAt,
			"url":          article.URL,
			"tweet_id":     article.TweetID,
			"tweet_id_str": article.TweetIDStr,
			"title":        article.Title,
			"summary":      article.Summary,
			"body":         article.Body,
			"tags":         article.Tags,
			"categories":   article.Categories,
			"publishedAt":  article.PublishedAt,
			"editedAt":     article.EditedAt,
			"keywords":     article.Keywords,
			"topics":       article.Topics,
			"entities":     article.Entities,
			"screen_name":  article.ScreenName,
		})

		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, fmt.Errorf("article %s record didn't create: %s", article.DocID, result.Err().Error())
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
	}

	if _, err := session.WriteTransaction(f); err != nil {
		return err
	}
	return nil
}

var publishedAtTxFunc = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:PUBLISHED_AT]->(b)
`

// CreatePublishedAtTxFunc Create article's similarity transaction function
func CreatePublishedAtTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(publishedAtTxFunc, map[string]interface{}{"source": source, "dest": dest})
	}
}

var authorOfTxFunc = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Article {uid: $dest})
	MERGE (a)-[:AUTHOR_OF]->(b)
`

// CreateAuthorOfTxFunc Create article's similarity transaction function
func CreateAuthorOfTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(authorOfTxFunc, map[string]interface{}{"source": source, "dest": dest})
	}
}

var writesForTxFunc = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:WRITES_FOR]->(b)
`

// CreateWritesForTxFunc Create article's similarity transaction function
func CreateWritesForTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(writesForTxFunc, map[string]interface{}{"source": source, "dest": dest})
	}
}

var hasEntityTxFunc = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Entity {uid: $dest})
	MERGE (a)-[:HAS_ENTITY]-(b)
`

// CreateHasEntityTxFunc Create article's similarity transaction function
func CreateHasEntityTxFunc(source string, dest string) neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		return tx.Run(hasEntityTxFunc, map[string]interface{}{"source": source, "dest": dest})
	}
}

// ListNodesWithQTxFunc Full-Text Search Query
var ListNodesWithQTxFunc = `
	CALL db.index.fulltext.queryNodes("text", '{{.q}}') YIELD node AS n
	WHERE n.crawledAt >= datetime("{{.from}}")
`

// ListNodesWithoutQTxFunc Without Q Query
var ListNodesWithoutQTxFunc = `
	MATCH (n:Article)--(f:Feed)
	WHERE n.crawledAt >= datetime("{{.from}}")
`

// ListNodesWithoutToTxFunc Without Q Query
var ListNodesWithoutToTxFunc = `
	AND n.crawledAt <= datetime("{{.to}}")
`

// ListCountNodesResultsTxFunc Return Total
var ListCountNodesResultsTxFunc = `
	{{if .withFullText}}
	CALL 
		db.index.fulltext.queryNodes("text", "{{.q}}") YIELD node AS n
	{{else}}
	MATCH 
		(n:Article)
	{{end}}

	WHERE 
		n.crawledAt >= datetime("{{.from}}")
	AND 
		n.crawledAt <= datetime("{{.to}}")

	{{if .withFeeds}}
	AND ({{.feeds}})
	{{end}}

	{{if .withURL}}
	AND n.url =~ "(?i){{.url}}"
	{{end}}

	{{if .withNotSimilar}}
	AND NOT 
		(n:Article)-[:SIMILAR_WITH]-(:Article)
	{{end}}

	{{if .withOptionalSimilar}}
	OPTIONAL MATCH 
		(n)-[r:SIMILAR_WITH]-(:Article)
	{{else if .withSimilar}}
	MATCH 
		(n:Article)-[r:SIMILAR_WITH]-(:Article)
	{{else}}
	{{end}}

	WITH 
		COUNT(DISTINCT n) AS total
	RETURN 
		total
`

// CountNodeTxFunc Count Single Node relationships
var CountNodeTxFunc = `
MATCH (n:Article { docId: $docId })-[r:SIMILAR_WITH]-(:Article)
RETURN count(r) as count
`

// GetNodeTxFunc Get Single Node data
var GetNodeTxFunc = `
	MATCH (n:Article { docId: $docId })

	WITH n
		MATCH (n)--(f:Feed)
		OPTIONAL MATCH (n)--(a:Author)
		OPTIONAL MATCH (n)-[r:SIMILAR_WITH]-(b:Article)

	WITH n, f, a, r, b
		OPTIONAL MATCH (b)--(fb:Feed)
		OPTIONAL MATCH (b)--(ab:Author)

	WITH n, f, a, b { .*, feed: fb { .* }, authors: collect(DISTINCT ab { .* }), score: r.score } as similar
	RETURN n { .*, feed: f { .* }, authors: collect(DISTINCT a { .* }), similar: collect(DISTINCT similar) } as json
`

// ListNodesTxFunc Nodes Query
var ListNodesTxFunc = `
	{{if .withFullText}}
	CALL 
		db.index.fulltext.queryNodes("text", "{{.q}}") YIELD node AS n
	{{else}}
	MATCH 
		(n:Article)
	{{end}}

	WHERE 
		n.crawledAt >= datetime("{{.from}}")
	AND 
		n.crawledAt <= datetime("{{.to}}")

	{{if .withFeeds}}
	AND ({{.feeds}})
	{{end}}

	{{if .withURL}}
	AND n.url =~ "(?i){{.url}}"
	{{end}}

	{{if .withNotSimilar}}
	AND NOT 
		(n:Article)-[:SIMILAR_WITH]-(:Article)
	{{end}}

	{{if .withOptionalSimilar}}
	OPTIONAL MATCH 
		(n)-[r:SIMILAR_WITH]-(:Article)
	{{else if .withSimilar}}
	MATCH 
		(n:Article)-[r:SIMILAR_WITH]-(:Article)
	{{else}}
	{{end}}

	WITH 
		COLLECT(DISTINCT n) AS articles
	UNWIND 
		articles AS n

	WITH n
		ORDER BY n.crawledAt DESC 
		SKIP {{.skip}} LIMIT {{.limit}}
	MATCH (n)--(f:Feed)

	{{if .withOptionalSimilar}}
	OPTIONAL MATCH 
		(n)-[r:SIMILAR_WITH]-(:Article)
	{{else if .withSimilar}}
	MATCH 
		(n:Article)-[r:SIMILAR_WITH]-(:Article)
	{{else}}
	WHERE NOT (n:Article)-[:SIMILAR_WITH]-(:Article)
	{{end}}

	WITH DISTINCT n { .*, feed: f { .* }{{.includeRels}}} AS data
	RETURN data
`
