# Ad Service API  

## Packages Used  

- **Fiber**: Used to implement REST API endpoints.  
- **pq**: PostgreSQL driver for Go, used for database connectivity and query execution.  
- **rand**: Generates a random selection of ads.  
- **godotenv**: Loads environment variables from a `.env` file.  

## 🔧 Environment Variables  

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

### **GET /ads**  
- Fetches a random ad from the available ones using the `rand` function.  

### **POST /ads/click**  
- Receives data sent by the client when a user clicks on an ad.  
- Records:  
  - `AdID`: Unique identifier for the clicked ad.  
  - `Timestamp`: Time when the click occurred.  
  - `VideoTime`: Video timestamp at which the ad was clicked.  
- Saves the click data in the PostgreSQL `clicks` table.  
