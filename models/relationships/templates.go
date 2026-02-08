package relationships

// drop duplicate nodes from neo4j with respect to memory
var dropDuplicates = `
	MATCH (a:Article)
	WITH a
	WITH a.uid AS uid, COLLECT(a) AS articles
	WHERE SIZE(articles) > 1
	
`
var indexes = `
	MATCH (a:Article)
	WITH a
	WITH a.uid AS uid, COLLECT(a) AS articles
	WHERE SIZE(articles) > 1

	UNWIND articles[1..] AS n
	DETACH DELETE n;
	
	MATCH (a:Article)
	WITH a.doc_id AS doc_id, COLLECT(a) AS articles
	WHERE SIZE(articles) > 1
	FOREACH (n IN TAIL(articles) | DETACH DELETE n);

	CREATE CONSTRAINT article_uid FOR (a:Article) REQUIRE a.uid IS UNIQUE;
	CREATE CONSTRAINT article_doc_id FOR (a:Article) REQUIRE a.doc_id IS UNIQUE;

	MATCH (f:Feed)
	WITH f.uid AS uid, COLLECT(f) AS feeds
	WHERE SIZE(feeds) > 1
	FOREACH (n IN TAIL(feeds) | DETACH DELETE n);

	MATCH (f:Feed)
	WITH f.feed_id AS feed_id, COLLECT(f) AS feeds
	WHERE SIZE(feeds) > 1
	FOREACH (n IN TAIL(feeds) | DETACH DELETE n);

	CREATE CONSTRAINT feed_uid FOR (f:Feed) REQUIRE f.uid IS UNIQUE;
	CREATE CONSTRAINT feed_id FOR (f:Feed) REQUIRE f.feed_id IS UNIQUE;
	CREATE CONSTRAINT feed_hostname FOR (f:Feed) REQUIRE f.hostname IS UNIQUE;

	MATCH (a:Author)
	WITH a.uid AS uid, COLLECT(a) AS authors
	WHERE SIZE(authors) > 1
	FOREACH (n IN TAIL(authors) | DETACH DELETE n);

	MATCH (a:Author)
	WITH a.author AS author, COLLECT(a) AS authors
	WHERE SIZE(authors) > 1
	FOREACH (n IN TAIL(authors) | DETACH DELETE n);

	CREATE CONSTRAINT author_uid FOR (a:Author) REQUIRE a.uid IS UNIQUE;
	CREATE CONSTRAINT author_name FOR (a:Author) REQUIRE a.author IS UNIQUE;

	MATCH (e:Entity)
	WITH e.uid AS uid, COLLECT(e) AS entities
	WHERE SIZE(entities) > 1
	FOREACH (n IN TAIL(entities) | DETACH DELETE n);

	MATCH (e:Entity)
	WITH e.entity_text AS uid, COLLECT(e) AS entities
	WHERE SIZE(entities) > 1
	FOREACH (n IN TAIL(entities) | DETACH DELETE n);

	CREATE CONSTRAINT entity_uid FOR (e:Entity) REQUIRE e.uid IS UNIQUE;
	CREATE CONSTRAINT entity_name FOR (e:Entity) REQUIRE e.entity_text IS UNIQUE;

	CREATE INDEX feed_meta FOR (f:Feed) ON (f.hostname, f.lang);
	CREATE INDEX article_meta FOR (a:Article) ON (a.hostname, a.url, a.lang);

	DROP CONSTRAINT article_tweetId IF EXISTS
	DROP CONSTRAINT article_tweetIdStr IF EXISTS
	DROP CONSTRAINT feed_twitter_id_str IF EXISTS
	DROP CONSTRAINT feed_twitter_id IF EXISTS

	CALL db.index.fulltext.createNodeIndex("text", ["Article"], ["title", "body", "keywords"]);

	
`

var nodeEntityTpl = `
	MERGE (n:Entity {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeOrganizationTpl = `
	MERGE (n:Organization {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeGPETpl = `
	MERGE (n:GPE {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodePersonTpl = `
	MERGE (n:Person {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeTopicTpl = `
	MERGE (n:Topic {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeAuthorTpl = `
	MERGE (n:Author {
		label: $label,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeFeedTpl = `
	MERGE (n:Feed {
		feed_id: $feed_id,
		name: $name, 
		hostname: $hostname,
		url: $url,
		type: $type
	})
	ON CREATE SET n.uid = $uid
	RETURN n.uid as uid
`
var nodeArticleTpl = `
	CREATE (n:Article {
		uid: $uid, 
		doc_id: $doc_id,
		lang: $lang,
		crawled_at: datetime($crawled_at),
		url: $url,
		title: $title,
		published_at: datetime($published_at),
		hostname: $hostname
	})
	RETURN n.uid as uid
`
var publishedAtTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:PUBLISHED_AT]->(b)
`
var authorOfTpl = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Article {uid: $dest})
	MERGE (a)-[:AUTHOR_OF]->(b)
`
var writesForTpl = `
	MATCH (a:Author {uid: $source})
	MATCH (b:Feed {uid: $dest})
	MERGE (a)-[:WRITES_FOR]->(b)
`

var hasEntityTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Entity {uid: $dest})
	MERGE (a)-[:HAS_ENTITY]-(b)
`
var topicTpl = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Topic {uid: $dest})
	MERGE (a)-[:IN_TOPIC]-(b)
`
var similarTxFunc = `
	MATCH (a:Article {uid: $source})
	MATCH (b:Article {uid: $dest})
	MERGE (a)-[:SIMILAR_WITH { score: $score }]->(b)
`

// CountSimilarTpl Count Single Node relationships
var CountSimilarTpl = `
	MATCH (a:Article { doc_id: $doc_id })-[r:SIMILAR_WITH]-(b:Article)
	RETURN count(r) as count
`

// GetNodeTpl Get Single Node data
var GetNodeTpl = `
	MATCH (n:Article { doc_id: $doc_id })

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

var GetSimilarTpl = `
	MATCH (n:Article { doc_id: $doc_id })-[r:SIMILAR_WITH]-(b:Article)	
	WITH r, b
	WITH b { .*, score: r.score } as node
	RETURN collect(DISTINCT node) as json
`
