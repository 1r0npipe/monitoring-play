# Play with monitoring solutions Prometheus

# Working:
1. to run use that command:
```bash
docker-compose up
```
it should download golang1.16-alpine and build image

2. Two Servers on port 8080 and 8081
2.1 to call port 8080, use that: "http://localhost:8080/data?age=78" (any namber after age, some check happens)
2.2 to use another server on port 8081: http://localhost:8081/year?year=1978 any year as digit (also some check happens)
3. Prometheus gathers the metrics as counter, gauge and histogramm
4. It is possible to access grafana as well on http://localhost:3000/ admin:admin by default
4.1 you need to set up at settings the source of data to see it IP address and port(http://localhost:9000), then save and test
5. you can access application metrics there: http://localhost:8080/metrics

Thank you