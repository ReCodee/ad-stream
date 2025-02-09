## Packages:

    - Fiber is used to implement REST Endpoints.
    - pq is used for connecting to PostgresSQL DB and executing the queries.
    - rand function is used send back random selection of Ad.
    - godotenv is used to load the environment variables.

## Environment Variables:

    - DB_PORT: Port at which db container is running.
    - DB_HOST: DB Container name, localhost for location connection.
    - DB_USER: DB Connection username.
    - DB_PASSWORD: DB Connection password.
    - DB_NAME: Name of the Ads DB.
    - CLICKS_TABLE: Clicks table name.
    - APP_PORT: Port at which app should run on. (Should be static for now i.e 8080)

### Endpoints:

    GET /ads
        - Fetches a random Ad among the available ones using the rand function.

    POST /ads/click
        - Receives the data which is send by the client application when the user clicks on a Ad.
        - It records Ad id (Unique identifier for the particular Ad), timestamp (Current time when the click happend) and VideoTime (video timestamp at which user clicked the Ad).
        - Saves the click data in PostgreSQL DB 'clicks' table.
