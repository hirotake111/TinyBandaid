### Environment variables

- `BACKEND_SERVERS`: comma-separated URLs for backend servers

### Example commands for development

```bash
# First of all get you Docker (Desktop) up and running (I always forget about it)
# Spin up dev environment
make up

# Make an HTTP request to the LB
curl http://localhost:3000

# Stop & remove all Docker containers
make down
```
