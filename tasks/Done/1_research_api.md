# Task: Understand the Sopra 'last-flight' API Endpoint

## Description
This task focuses on understanding and integrating with the existing `sopra` API, specifically the `last-flight` endpoint, as specified by the user.

## API Details
The user has confirmed that the `sopra` API provides a `last-flight` endpoint which directly returns the necessary formatted data:
*   IATA City <Airplane Icon> IATA City
*   Time of seeing the flight
*   Airline

This means the ESP device will directly consume this endpoint's output without requiring complex external flight API integrations, data transformations, or lookups for IATA cities/airlines.

## Next Steps:
1.  **Identify the exact URL** of the `sopra` `last-flight` endpoint.
2.  **Determine the request method** (GET, POST) and any required parameters (e.g., geographical coordinates, user ID, authentication tokens).
3.  **Understand the exact JSON response structure** from this endpoint to correctly parse the "IATA City <Airplane Icon> IATA City", "time of seeing the flight", and "Airline" fields.
4.  If authentication is required, identify the authentication mechanism (e.g., API key, token, basic auth).

## Status: Completed
The research for the `last-flight` API endpoint has been completed. The results are documented in `1_research_api_results.md`.
