# Rate Limit by Status Code

MVP for a Gorilla Mux middleware which limits the rate of unauthorized 
responses to a HTTP client.

Designed to be deployed to [fly.io](https://fly.io) with a Redis-backed counter,
but `pkg/limit` should be usable anywhere.

Clients are identified by their IP address. Something smarter than this might be
required in the future as Apple's privacy proxy gains wider use.

## Environment Variables (Default Value)

All environment variables are read in 

- PORT (8080) - TCP port to listen on
- LOG_LEVEL (4) - logrus level, defaults to info, higher is more-granular
- RATE_COUNT (3) - number of unauthorized attempts before blocking client
- RATE_TTL (10) - time window for rate limit in seconds
- FLY_REDIS_CACHE_URL - redis:// url to use for persistent cache
