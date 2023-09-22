# Usage
```bash
docker build -t stddev-service .
docker run -p 80:80 stddev-service
```

```bash
curl "http://localhost/random/stddevs?requests=2&length=5"
```
