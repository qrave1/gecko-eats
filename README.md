### cron command

```bash
docker run --rm \
  -v ./config.yaml:/app/config.yaml \
  docker.io/qrave1/gecko-eats:latest notify -c /app/config.yaml
```