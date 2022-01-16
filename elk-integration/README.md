# Play with monitoring solutions Prometheus

# Working:
1. to run use that command:
```bash
docker-compose up
```
it should download golang1.16-alpine and build image
make sure the permissions for elk-integration/elasticsearh is 777

2. Two Servers on port 8080 and 8081
2.1 to call port 8080, use that: "http://localhost:8080/data?age=78" (any namber after age, some check happens)
2.2 to use another server on port 8081: http://localhost:8081/year?year=1978 any year as digit (also some check happens)
3. ElasticSearch+Kibana+FluentD are run withing same configuration, only you may need to create index and Dashboard

p.s. I played with Kibana (it was my first time, it quite sophisticated software and I decided review it just quickly) it is not so easy
and this is enough to make the homework done

Thank you