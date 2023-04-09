# Rating Orama

Rating Orama is a web application for displaying TV show ratings and statistics. It is composed of 3 main parts:

1. **Core**: Written in Go and Fiber, responsible for orchestrating everything and displaying the data using a template engine.
2. **Harvester**: Written in Python, Flask, and Cinemagoer, responsible for collecting data for the core.
3. **Database**: PostgreSQL for storing data.

## Running the project

There are two ways to run the project: launching each part individually or building the Dockerfile and running it using Docker Compose. Here's an example of the `docker-compose.yml` file for the latter option:

```yaml
version: '3'

services:
  harvester:
    container_name: harvester-ratingorama
    image: harvester:0.1.0
    networks:
      - ratingorama
  core:
    container_name: core-ratingorama
    image: core:0.1.0
    environment:
      DATASOURCE: ${DATASOURCE}
      HARVESTER_API: ${HARVESTER_API}
      IS_PRODUCTION: ${IS_PRODUCTION}
    ports:
      - "3000:3000"
    networks:
      - ratingorama
  db:
    container_name: db-ratingorama
    image: postgres:15.2-alpine
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5433:5432"
    volumes:
      - ./schema.sql:/docker-entrypoint-initdb.d/schema.sql
      - ./data:/var/lib/postgresql/data
    networks:
      - ratingorama

networks:
  ratingorama:
```

## Contributions

If you have ideas for improvements or bug fixes, feel free to contribute! To do so, simply clone the repository, create a new branch, and submit a pull request.
