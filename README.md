# Rating Orama

Rating Orama is a web application for displaying TV show ratings and statistics. 
It is composed of 2 main parts:

1. **Core**: Written in Go and Fiber, responsible for orchestrating everything 
and displaying the data using a template engine.
3. **Database**: PostgreSQL for storing data.

## Running the project

There are two ways to run the project: launching each part individually or 
building the Dockerfile and running it using Docker Compose. Here's an example 
of the `docker-compose.yml` file for the latter option:

```yaml
version: '3'

services:
  core:
    container_name: core-ratingorama
    image: core:latest
    environment:
      DATASOURCE: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@db:5432/${POSTGRES_DB}?sslmode=disable
    ports:
      - "8080:8080"
    networks:
      - ratingorama
  db:
    container_name: db-ratingorama
    image: postgres:16.3-alpine3.20
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - rating-orama_data:/var/lib/postgresql/data
    networks:
      - ratingorama

networks:
  ratingorama:

volumes:
  rating-orama_data:
```

## Contributions

If you have ideas for improvements or bug fixes, feel free to contribute! To do 
so, simply clone the repository, create a new branch, and submit a pull request.
