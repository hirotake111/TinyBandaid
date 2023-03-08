## 1. Simple Round Robin loadbalancing

```bash
# First of all get your Docker (Desktop) up and running (as I always forget about it)
# Spin up dev environment
make up

# Open another terminal and make an HTTP request to the LB
curl http://localhost:3000

# Stop & remove all Docker containers
make down
```
