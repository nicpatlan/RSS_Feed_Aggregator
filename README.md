# RSS_Feed_Aggregator
A project, written in Go, for aggregating RSS feeds for users. Utilizes a 
PostgreSQL database using Goose for database migration and sqlc to generate 
Go code for SQL queries in place of an ORM.

Batches of database feeds will be aggregated in groups of 10 once per minute
with preference to the least recently fetched feeds first. This can be adjusted 
in main.go based on preference.

## Installation & Setup
Inside a go module:
>     go get github.com/nicpatlan/RSS_Feed_Aggregator

PostgreSQL can be installed using one of the methods below.
MacOS with Homebrew:
>     brew install postgresql@15

To start the PostgreSQL service run:
>     brew services start postgresql

To enter the psql shell run:
>     psql postgres

Linux:
>     sudo apt update
>     sudo apt install postgresql postgresql-contrib

Linux users may need to set a password as an additional step.
>     sudo passwd postgres

To start PostgreSQL service run:
>     sudo service postgresql start

To enter the psql shell run:
>     sudo -u postgres psql

Create a database using SQL syntax with the following table names:
- users
- feeds
- users_feeds
- posts

Install the latest versions of the needed CLI tools.
Goose:
>     go install github.com/pressly/goose/v3/cmd/goose@latest

Connection string will be of the format:
>     protocol://username:password@host:port/database

Run the command below for all up migrations in the project using the 
connection string above replaced with the appropriate username,
password (if needed), host, port and database name.
>     goose postgres connection_string up

Down migration can be done with the command:
>     goose postgres connection_string down

Other commands for [postgreSQL](https://www.postgresql.org/docs/current/app-psql.html#:~:text=psql%20is%20a%20terminal%2Dbased,or%20from%20command%20line%20arguments.)

PostgreSQL driver can be installed by running:
>     go get github.com/lib/pq

sqlc can be installed if needed to make adjustments to queries.
>     go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## API
The server runs on the localhost port 8080 but can be adjusted to another web
server in main.go if desired.

Typical access pattern will be http://localhost:8080 followed by an endpoint 
detailed below.

### Server Metric Endpoints
- GET /v1/healthz  "returns the server status"
- GET /v1/err      "error formatting"

### User Endpoints
- POST /v1/users  "creates a new user with name given in body"
- GET /v1/users   "authenticated endpoint that returns user information"

The request body to the endpoint for creating a user should conform to the example 
below:
> {
>    "name": "user_name"
> }

The endpoint for retrieving user information should include in the header:
> {
    "Authorization": "ApiKey user_api_key"
> }

### Feed Endpoints
- POST /v1/feeds  "creates a feed with the provided name and url"
- GET /v1/feeds   "returns all feeds currently stored in the database"

The request body to the endpoint for creating a feed should conform to the example
below:
> {
    "name": "name_of_feed",
    "url":  "url_of_feed"
> }

### Feed Follow Endpoints
- POST /v1/feed_follows                   "sets a user as follower of a feed"
- GET /v1/feed_follows                    "retrieves all feeds followed by a user"
- DELETE /v1/feed_follows/{feedFollowID}  "removes the feed follow with the provided ID"

Setting a user as a follower of a feed and retrieving feeds followed by a user are 
authenticated endpoints and require a header including:
> {
    "Authorization": "ApiKey user_api_key"
> }

Deleting a feed follow requires the feed follow id which will be returned when using the 
authenticated GET /v1/feed_follows endpoint.

### Feed Posts Endpoint
- GET /v1/posts  "retrieves the posts stored for the users followed feeds"

This is an authenticated endpoint and requires a header including:
> {
    "Authorization": "ApiKey user_api_key"
> }

There is an optional query parameter that limits the number of posts returned in the 
response that should follow the form below. The returned posts will be in descending 
order from most recent onward. If no query parameter is given the default is 5.
> ?limit=number_of_posts_desired