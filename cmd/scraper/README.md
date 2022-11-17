# MediaWatch Scraper Service

gRPC micro-service for ascraping articles. Make sure `mongodb` is up and running to retrieve `trim passages`.

## Docker Container

```bash
# Build
make docker
# Run
docker run --rm \
	-p 500050:50050 -e MONGODB_URL mymongo \
	reg.plagiari.sm/psm-svc-scraper:1.1.0
```

### Enviroment Variables

```bash
# Default Service gRPC
SERVER_ADDRESS=localhost:50050
# Default proto path
PROTO_PATH=/app/proto
# Default MongoDB Url
MONGODB_URL=mongodb://localhost:27017/
MONGODB_DB=mediawatch
```
