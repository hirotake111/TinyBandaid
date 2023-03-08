## Simple Round Robin loadbalancing

### Example commands for development

```bash
# First of all get you Docker (Desktop) up and running (I always forget about it)
# Spin up dev environment
make up

# Open another terminal and make an HTTP request to the LB
curl http://localhost:3000

# Stop & remove all Docker containers
make down
```
