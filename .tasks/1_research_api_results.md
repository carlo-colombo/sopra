# Sopra 'last-flight' API Endpoint Research Results

This document summarizes the findings regarding the `sopra` `last-flight` API endpoint, based on the analysis of `server/server.go`.

## 1. Exact URL

The endpoint is accessible via a `GET` request at:
`/last-flight`

The full URL will depend on the server's address and configured port (e.g., `http://localhost:8080/last-flight`).

## 2. Request Method and Parameters

*   **Method:** `GET`
*   **Parameters:** This endpoint does not require any query parameters or request body.

## 3. Exact JSON Response Structure

The endpoint returns a JSON object with the following structure:

```json
{
    "flight": "string",             // Flight identifier (e.g., "AAL123")
    "operator": "string",           // Short name of the airline operator (e.g., "American Airlines")
    "destination_city": "string",   // Destination city name
    "destination_code_iata": "string", // Destination IATA code (e.g., "LAX")
    "destination_code_icao": "string", // Destination ICAO code
    "source_city": "string",        // Source city name
    "source_code_iata": "string",   // Source IATA code (e.g., "JFK")
    "source_code_icao": "string",   // Source ICAO code
    "last_time_seen": "string",     // Timestamp of when the flight was last seen (ISO 8601 format, e.g., "2006-01-02T15:04:05Z")
    "airplane_model": "string"      // Type of aircraft (e.g., "A320")
}
```

**Mapping to User Requirements:**
*   **"IATA City <Airplane Icon> IATA City"**: Can be constructed from `source_code_iata` and `destination_code_iata` (e.g., "JFK ✈️ LAX"). The airplane icon should be added by the client.
*   **"Time of seeing the flight"**: Corresponds to the `last_time_seen` field.
*   **"Airline"**: Corresponds to the `operator` field.

## 4. Authentication Mechanism

This endpoint does **not** require any authentication. It is publicly accessible.
