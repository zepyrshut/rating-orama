# Changelog

All notable change notes in this project will be documented in this file.

## [0.1.0-2] - 2023-04-09
### Changed
- `docker-compose.yml`: Changed the volume path for the `init.sql` file:
  - From: `./schema.sql:/docker-entrypoint-initdb.d/schema.sql`
  - To: `./init.sql:/docker-entrypoint-initdb.d/init.sql`

### Removed
- `docker/schema.sql`: The file has been removed. The schema configuration has been moved to the `init.sql` file.
