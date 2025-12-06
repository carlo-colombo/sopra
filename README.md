Sopra is a Go service that identifies airplanes flying over a specific location. It uses the OpenSky Network API to fetch real-time flight data and provides a simple REST API to access this information. The service can be run as a standalone application or as a Docker container.

## Configuration

The application can be configured using a `config.yml` file, environment variables, or command-line flags. The order of precedence is as follows:

1.  Command-line flags
2.  Environment variables
3.  `config.yml` file
4.  Default values

### `config.yml`

Create a `config.yml` file in the root of the project to set your configuration options. Here's an example:

```yaml
port: 8080
db_path: "sopra.db"

opensky_client:
  id: "your_opensky_client_id"
  secret: "your_opensky_client_secret"

flightaware:
  api_key: "your_flightaware_api_key"

service:
  latitude: 47.3769
  longitude: 8.5417
  radius: 100.0
```

### Environment Variables

You can also set configuration options using environment variables. Here's a list of the available environment variables:

| Variable                | Description                               |
| ----------------------- | ----------------------------------------- |
| `PORT`                  | The port to run the server on.            |
| `DB_PATH`               | The path to the SQLite database file.     |
| `OPENSKY_CLIENT_ID`     | Your OpenSky API client ID.               |
| `OPENSKY_CLIENT_SECRET` | Your OpenSky API client secret.           |
| `FLIGHTAWARE_API_KEY`   | Your FlightAware API key.                 |
| `DEFAULT_LATITUDE`      | The default latitude for flight searches. |
| `DEFAULT_LONGITUDE`     | The default longitude for flight searches.|
| `DEFAULT_RADIUS`        | The default radius for flight searches.   |
| `WATCH`                 | Enable watch mode.                        |
| `WATCH_INTERVAL`        | The interval to watch for flights in seconds. |

### Command-line Flags

Finally, you can use command-line flags to set certain configuration options:

| Flag       | Description                                  |
| ---------- | -------------------------------------------- |
| `--print`  | Print the result and logs to stdout.         |
| `--watch`  | Watch for flights and log them.              |
| `--interval`| The interval to watch for flights in seconds.|

## Local Execution

1.  **Set Environment Variables:**
    Ensure you have your OpenSky API `client_id` and `client_secret` set as environment variables.

    ```bash
    export OPENSKY_CLIENT_ID="your_client_id"
    export OPENSKY_CLIENT_SECRET="your_client_secret"
    ```

2.  **Build the Application:**
    ```bash
    go build
    ```

3.  **Run the Application:**
    ```bash
    ./sopra
    ```

## API Endpoints

The application exposes the following API endpoints:

### `/flights`

Returns a list of all flights currently in the specified radius.

**Example Response:**

```json
[
  {
    "icao24": "a8b4c2",
    "callsign": "SWR123",
    "origin_country": "Switzerland",
    "time_position": 1678886400,
    "last_contact": 1678886400,
    "longitude": 8.5417,
    "latitude": 47.3769,
    "baro_altitude": 10000,
    "on_ground": false,
    "velocity": 250,
    "true_track": 180,
    "vertical_rate": 0,
    "sensors": null,
    "geo_altitude": 10000,
    "squawk": "1234",
    "spi": false,
    "position_source": 0
  }
]
```

### `/last-flight`

Returns the last flight that was recorded in the database.

**Example Response:**

```json
{
  "flight": "SWR123",
  "operator": "Swiss",
  "destination_city": "Zurich",
  "destination_code": "ZRH",
  "source_city": "Geneva",
  "source_code": "GVA",
  "last_time_seen": "2023-03-15T12:00:00Z",
  "airplane_model": "A320"
}
```

### `/all-flights`

Returns all flights that have been recorded in the database.

**Example Response:**

```json
[
  {
    "flight": "SWR123",
    "operator": "Swiss",
    "destination_city": "Zurich",
    "destination_code": "ZRH",
    "source_city": "Geneva",
    "source_code": "GVA",
    "last_time_seen": "2023-03-15T12:00:00Z",
    "airplane_model": "A320"
  }
]
```

## Project Structure

The project is organized into the following directories:

| Directory  | Description                               |
| ---------- | ----------------------------------------- |
| `client`   | Contains the OpenSky and FlightAware API clients. |
| `config`   | Handles application configuration.        |
| `database` | Manages the SQLite database.              |
| `haversine`| Provides functions for calculating distances between coordinates. |
| `model`    | Defines the data models for the application. |
| `server`   | Contains the HTTP server and API endpoints. |
| `service`  | Implements the core business logic.      |
