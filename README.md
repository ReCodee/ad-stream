# Ad Service API  

## Packages Used  

- **Fiber**: Used to implement REST API endpoints.  
- **pq**: PostgreSQL driver for Go, used for database connectivity and query execution.  
- **rand**: Generates a random selection of ads.  
- **godotenv**: Loads environment variables from a `.env` file.  
- **gofiber/websocket/v2**: Used to implement WebSocket connections for real-time ad updates.  

## Environment Variables  

| Variable      | Description |
|--------------|-------------|
| `DB_PORT`    | Port where the database container is running. |
| `DB_HOST`    | Database container name (use `localhost` for local connection). |
| `DB_USER`    | Database connection username. |
| `DB_PASSWORD`| Database connection password. |
| `DB_NAME`    | Name of the Ads database. |
| `CLICKS_TABLE` | Name of the table storing ad clicks. |
| `APP_PORT`   | Port where the app runs (should be static, e.g., `8080`). |

## API Endpoints  

### **GET /ads (Replaced by WebSocket)**  
- Fetches a random ad from the available ones using the `rand` function.  

### **POST /ads/click**  
- Receives data sent by the client when a user clicks on an ad.  
- Records:  
  - `AdID`: Unique identifier for the clicked ad.  
  - `Timestamp`: Time when the click occurred.  
  - `VideoTime`: Video timestamp at which the ad was clicked.  
  - `Position`: Position of the ad on the screen when clicked.  
  - `HoverTime`: Time spent hovering over the ad before clicking.  
- Saves the click data in the PostgreSQL `clicks` table.  

### **GET /ws**  
- Establishes a WebSocket connection with the client.  
- Sends a random ad to the client every 10 seconds.  


### Running the server

```bash
go run main.go
```

### Production Build with Docker

```bash
docker network create ad-service-network
docker build -t ad-service .

docker run -d --name postgres --network ad-service-network ^
-e POSTGRES_USER=username ^
-e POSTGRES_PASSWORD=password ^
-e POSTGRES_DB=ads ^
postgres:latest

docker run -d --name ad-service --network ad-service-network -p 8080:8080 ^
-e DB_HOST=postgres ^
-e DB_USER=username ^
-e DB_PASSWORD=password ^
-e DB_NAME=ads ^
-e DB_PORT=5432 ^
ad-service
```


