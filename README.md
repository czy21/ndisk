# ndisk

```bash
bash build.sh && docker build --tag registry.cluster.com/library/ndisk:master --no-cache --file Dockerfile .
docker run --rm --env CONFIG_FILE=/app/app.yaml --name ndisk --mount type=bind,src=/home/opsor/ndisk/backend/build/app.yaml,dst=/app/app.yaml -p 8080:8080 registry.cluster.com/library/ndisk:master
```