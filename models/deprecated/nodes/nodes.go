package nodes

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
