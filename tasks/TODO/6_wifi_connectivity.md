# Task 6: Wi-Fi Connectivity with MicroPython on ESP32

## Description
This task focuses on establishing Wi-Fi connectivity for the ESP32 using MicroPython. This is a fundamental step for any network-enabled IoT application, allowing the device to communicate with external services (like the Go server).

## Prerequisites
*   MicroPython firmware flashed to ESP32.
*   A Wi-Fi network (SSID and password) available for connection.

## Steps

### 6.1. Write MicroPython Code for Wi-Fi Connection
1.  **Import necessary packages:** You will typically need the `network` module for Wi-Fi control.
2.  **Configure Wi-Fi credentials:** Embed your Wi-Fi SSID and password directly in the code for initial testing, or implement a more secure configuration method later (e.g., from a configuration file or EEPROM).
3.  **Implement connection logic:**
    *   Initialize the `WLAN` interface for station mode (e.g., `wlan = network.WLAN(network.STA)`).
    *   Activate the interface (`wlan.active(True)`).
    *   Connect to the specified access point using the SSID and password (`wlan.connect('SSID', 'PASSWORD')`).
    *   Implement retry logic and status checks (e.g., `while not wlan.isconnected(): time.sleep(1)`) to ensure a robust connection.
    *   Print status messages (e.g., IP address `wlan.ifconfig()`) to the serial console to aid in debugging.

### 6.2. Test Wi-Fi Connectivity
1.  **Simple Network Request:** After successfully connecting to Wi-Fi, perform a simple network operation to confirm connectivity. Examples include:
    *   Making an HTTP GET request to a public API (e.g., `urequests` module to `http://httpbin.org/ip` to get the device's external IP).
    *   Pinging a well-known host (e.g., `usocket` or `uping` if available).
    *   Connecting to a TCP server (`usocket`).
2.  **Print Results:** Output the results of the network test to the serial console.

### 6.3. Upload and Run
1.  **Connect to ESP32:** Use `ampy` or `webrepl` to connect to your MicroPython device.
2.  **Upload the Python script:** Upload your `main.py` (or other script) to the ESP32 filesystem.
    ```bash
    ampy -p /dev/ttyUSB0 put main.py
    ```
    (Replace `/dev/ttyUSB0` with your ESP's serial port.)

## Verification
*   **Serial Console Output:** Observe the serial output for messages indicating successful Wi-Fi connection (e.g., assigned IP address) and the results of your network test (e.g., HTTP response, ping success).
*   **Troubleshooting:** If connection fails, double-check SSID/password, ensure the ESP32 is within Wi-Fi range, and review MicroPython's `network` module documentation for common issues. Use the MicroPython REPL for interactive debugging.

## Next Steps
Once reliable Wi-Fi connectivity is established, the ESP32 can begin communicating with the Go server or other network resources, allowing for data exchange and remote control.
