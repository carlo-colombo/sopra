# Task 5: Display "Hello World" with MicroPython on ESP32

## Description
This task involves implementing a basic "Hello World" text display on the ST7789/ST7735 screen using MicroPython on the ESP32. This will confirm the display hardware integration and basic MicroPython code execution on the device.

## Prerequisites
*   MicroPython firmware flashed to ESP32.
*   Display (ST7789/ST7735) connected to the ESP32, typically via SPI. Ensure correct pin assignments for SPI (SCK, MOSI, MISO, CS, DC, RST) and backlight control.

## Steps

### 5.1. Identify/Adapt Display Driver
1.  **Search for MicroPython ST7789/ST7735 drivers:** Look for existing MicroPython libraries (often `.py` files) that support these display controllers. Popular options might be found in the MicroPython community forums, GitHub repositories (e.g., `micropython-lib`), or specific module documentation.
2.  **Review driver usage:** Understand how to initialize the SPI bus and the display controller using the chosen driver.
3.  **Adapt if necessary:** You may need to modify a driver for your specific pinout or display configuration.

### 5.2. Write MicroPython Code for Display
1.  **Initialize SPI:** Configure the appropriate ESP32 pins for SPI communication using the `machine` module (e.g., `machine.SPI`).
2.  **Initialize Display:** Import and instantiate the display driver, passing the SPI object and pin assignments (e.g., `st7789.ST7789(...)`).
3.  **Clear Screen:** Clear the display to a solid color using the driver's methods.
4.  **Render Text:** Use the driver's API to draw "Hello World" on the screen. This might involve setting font, color, and position (e.g., using `framebuf` or specific text methods).

### 5.3. Upload and Run
1.  **Connect to ESP32:** Use `ampy` or `webrepl` to connect to your MicroPython device.
2.  **Upload the Python script:** Upload your `main.py` (or other script) to the ESP32 filesystem.
    ```bash
    ampy -p /dev/ttyUSB0 put main.py
    ```
    (Replace `/dev/ttyUSB0` with your ESP's serial port, e.g., `COMx` on Windows.)
3.  **Soft reset/Reboot:** The script should run automatically on boot. If it doesn't, perform a soft reset or reboot the ESP32.

## Verification
*   **Observe Display:** After uploading and resetting, the ESP32 should restart, and "Hello World" should appear on the connected ST7789/ST7735 display.
*   **Troubleshooting:** If the display doesn't work, double-check wiring, pin assignments, SPI mode, and the display initialization sequence in your code. Ensure the MicroPython driver is correctly configured for your specific display variant (e.g., resolution, rotation). Use the MicroPython REPL for debugging.

## Next Steps
Once "Hello World" is successfully displayed, proceed to Task 6 for implementing Wi-Fi connectivity.
