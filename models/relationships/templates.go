package relationships

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
