sopra is a golang service that uses openskky rest api (https://openskynetnetwork.github.io/opensky-api/rest.html) to identify what airplane is over the sky in this moment in a pprecise location.

the service runs in the background and it should authenticate using the opensky api in a way that is doable vheadless without a browser in a server.

It returns a json that contains the list of planes in a configurable radius.

## OpenSky API Authentication Setup

To use this service, you need to authenticate with the OpenSky Network API. This service uses the OAuth2 Client Credentials Flow, which requires a Client ID and Client Secret. Follow these steps to obtain them:

1.  **Create an API Client:**
    *   Log in to your OpenSky account on their website (https://opensky-network.org).
    *   Navigate to your Account page.
    *   Create a new API client to obtain your `client_id` and `client_secret`.

2.  **Set Environment Variables:**
    Once you have your `client_id` and `client_secret`, set them as environment variables:
    ```bash
    export OPENREDISKY_CLIENT_ID="your_client_id"
    export OPENREDISKY_CLIENT_SECRET="your_client_secret"
    ```
    Replace `"your_client_id"` and `"your_client_secret"` with your actual credentials.

## Build and Run

1.  **Set Environment Variables:**
    Ensure you have your OpenSky API `client_id` and `client_secret` set as environment variables `OPENREDISKY_CLIENT_ID` and `OPENREDISKY_CLIENT_SECRET`. You can find instructions on how to get them in the "OpenSky API Authentication Setup" section above.

    ```bash
    export OPENREDISKY_CLIENT_ID="your_client_id"
    export OPENREDISKY_CLIENT_SECRET="your_client_secret"
    ```

2.  **Build the Docker Image:**
    Navigate to the project's root directory and run:
    ```bash
    docker build -t sopra .
    ```

3.  **Run the Docker Container:**
    ```bash
    docker run -p 8080:8080 -e OPENREDISKY_CLIENT_ID=$OPENREDISKY_CLIENT_ID -e OPENREDISKY_CLIENT_SECRET=$OPENREDISKY_CLIENT_SECRET sopra
    ```
    This will start the server, listening on port `8080`.

4.  **Access the API:**
    You can access the API endpoint at `http://localhost:8080/flights` to get the list of flights.
