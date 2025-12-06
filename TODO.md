# TODO - OpenSky to FlightAware Integration

This document outlines the tasks required to integrate OpenSky and FlightAware to display flights over a specific location.

## 1. Configuration

-   [x] Create a `config` package to hold all configuration settings.
-   [x] Add a `config.go` file to manage API keys, secrets, and other settings from environment variables or a configuration file.
-   [x] Define a struct to hold configuration for:
    -   OpenSky API (client ID, client secret, base URL, token URL)
    -   FlightAware AeroAPI (API key, base URL)
    -   Location settings (latitude, longitude, radius in km).

## 2. OpenSky Client

-   [x] Create a `client` package for API clients.
-   [x] Implement `opensky.go` to communicate with the OpenSky API.
-   [x] Implement OAuth2 client credentials flow to obtain an access token.
    -   The token should be cached and refreshed when it expires.
-   [x] Implement a function to fetch flight state vectors within a given bounding box.
    -   The bounding box should be calculated based on the configured latitude, longitude, and radius.
-   [x] Parse the JSON response into a `StateVector` model.
    -   The `model` package should contain the `StateVector` struct with fields like `icao24`, `callsign`, `lat`, `lon`, etc.
-   [ ] Filter the returned flights to be within the specified radius from the center point using the Haversine formula.

## 3. FlightAware Client

-   [x] Implement `flightaware.go` in the `client` package.
-   [x] Implement a function to fetch flight information by its identifier (callsign) from the AeroAPI.
    -   The request must be authenticated with the AeroAPI key.
-   [ ] Parse the JSON response into a `FlightInfo` model.
    -   The `model` package should contain the `FlightInfo` struct with fields for `ident`, `operator`, `aircraft_type`, `origin`, and `destination`.

## 4. Service Layer

-   [x] Create a `service` package to orchestrate the integration.
-   [x] Implement a `service.go` file with a `FlightService`.
-   [x] The `FlightService` should:
    -   [x] Fetch state vectors from the `OpenSkyClient`.
    -   [x] For each flight with a valid callsign, enrich the data by calling the `FlightAwareClient`.
    -   [x] Return a list of enriched `FlightInfo` objects.

## 5. Main Application

-   [x] Update `main.go` to tie everything together.
-   [x] Initialize the configuration, clients, and the FlightService.
-   [x] Create a simple command-line interface or a web server to trigger the flight fetching process.
-   [x] The application should periodically fetch and display the flight information (flight number, origin, destination, etc.).
-   [x] Implement logging to provide insights into the application's behavior.
-   [x] Enhance `/last-flight` endpoint to include `last_time_seen` and `airplane_model` in the JSON response.

## 6. Models

-   [x] Create a `model` package for data structures.
-   [x] Define `opensky.go` with structs for `StateVector`.
-   [ ] Define `flightaware.go` with structs for `FlightInfo`, `AirportInfo`, etc.

## 7. Geolocation Utilities

-   [x] Create a `haversine` package for geolocation calculations.
-   [ ] Implement `haversine.go` with functions to:
    -   Calculate the distance between two geographical points (Haversine formula).
    -   Calculate the bearing between two points.
-   [x] Calculate a bounding box around a center point given a radius.