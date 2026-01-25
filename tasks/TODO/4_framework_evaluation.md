# Task: Framework Evaluation and Initial Setup

## Description
This task aims to evaluate different programming frameworks for the ESP32/ESP8266 by performing initial setup and a basic "Hello World" display test for each. This will help you choose the framework that best fits your preferences and project requirements.

## Frameworks to Evaluate:
1.  **MicroPython** (Python-based, good for rapid prototyping, user's preference)
2.  **Go Embed** (TinyGo on ESP-IDF, Go language)
3.  **JS Stuff** (e.g., Espruino, NodeMCU firmware; JavaScript)
4.  **AtomVM / Elixir** (Erlang VM for embedded, Elixir language)

## General Prerequisites for All Frameworks:
*   USB-to-Serial drivers installed (if using bare ESP modules).
*   ESP Flash Download Tool or `esptool.py` for flashing firmware.
*   Python 3 installed (often needed for `esptool.py` and other utilities).

---

## 1. MicroPython Evaluation

### 1.1. Flash MicroPython Firmware
*   **Download:** Get the latest MicroPython firmware for your ESP32 or ESP8266 from [micropython.org/download](http://micropython.org/download). Choose the generic firmware.
*   **Erase Flash:**
    ```bash
    esptool.py --port /dev/ttyUSB0 erase_flash
    ```
    (Replace `/dev/ttyUSB0` with your ESP's serial port, e.g., `COMx` on Windows).
*   **Flash Firmware:**
    ```bash
    esptool.py --port /dev/ttyUSB0 --baud 460800 write_flash --flash_size=detect 0x1000 <path_to_firmware.bin>
    ```
    (Adjust baud rate if flashing fails, `460800` is usually reliable).
*   **Connect:** Use a serial terminal (e.g., `screen`, `PuTTY`, `minicom`, `Thonny IDE`) to connect to the ESP and access the MicroPython REPL.

### 1.2. Basic Display "Hello World" (MicroPython)
*   You will need MicroPython drivers for your ST7789/ST7735 display. Libraries like `micropython-ili9341` or `micropython-st7789` (often compatible with ST7735) are available.
*   **Steps:**
    1.  Install `upip` on your ESP (if not already present): `import upip; upip.install('micropython-lib')`
    2.  Find and install an appropriate display driver using `upip` or manually upload the driver files to the ESP's filesystem.
    3.  Write a small `main.py` script to initialize the SPI bus, the display driver, and draw some text like "Hello World" on the screen.
    4.  Upload `main.py` to your ESP.

### 1.3. Wi-Fi Connectivity (MicroPython)
*   **Steps:**
    1.  Create a `boot.py` file with your Wi-Fi credentials.
    2.  Use the `network` module to connect to your Wi-Fi network.
    3.  Test connectivity by pinging an external server.

---

## 2. Go Embed (TinyGo on ESP-IDF) Evaluation

### 2.1. Setup ESP-IDF
*   Go embed on ESP32 typically uses TinyGo, which compiles Go code to WebAssembly or native code and integrates with the ESP-IDF toolchain. This is a more involved setup.
*   **Install ESP-IDF:** Follow the official Espressif instructions to install ESP-IDF for your OS. This includes toolchains, Msys/MSYS2 for Windows, etc.
*   **Install TinyGo:** Follow TinyGo's instructions for installing TinyGo and setting up the ESP32 target.

### 2.2. Basic Display "Hello World" (Go Embed)
*   TinyGo has experimental support for graphics and display drivers.
*   **Steps:**
    1.  Find or adapt an existing TinyGo driver for ST7789/ST7735.
    2.  Write a simple Go program to initialize SPI and the display, then render text.
    3.  Compile with TinyGo and flash using `go run . -target esp32 -port /dev/ttyUSB0` (or similar).

### 2.3. Wi-Fi Connectivity (Go Embed)
*   TinyGo provides network capabilities for ESP32.
*   **Steps:**
    1.  Write Go code to connect to Wi-Fi using TinyGo's network packages.
    2.  Test HTTP requests.

---

## 3. JS Stuff (Espruino / NodeMCU) Evaluation

### 3.1. Flash Firmware (Espruino)
*   **Download:** Get the Espruino firmware for your ESP32 or ESP8266 from [espruino.com/Download](http://www.espruino.com/Download).
*   **Flash:** Use the Espruino Web IDE or `esptool.py` to flash the firmware.

### 3.2. Basic Display "Hello World" (JS)
*   Espruino has built-in or readily available modules for many displays.
*   **Steps:**
    1.  Connect to the Espruino Web IDE.
    2.  Initialize SPI and the ST7789/ST7735 module.
    3.  Use graphics commands to draw text and shapes.

### 3.3. Wi-Fi Connectivity (JS)
*   Espruino's `require("Wifi")` module handles Wi-Fi.
*   **Steps:**
    1.  Connect to Wi-Fi within your JavaScript code.
    2.  Test HTTP requests using `require("http")`.

---

## 4. AtomVM / Elixir Evaluation

### 4.1. Setup AtomVM / Elixir
*   AtomVM brings the Erlang VM to microcontrollers, allowing Elixir to run. This is a more advanced setup.
*   **Install Erlang/Elixir:** You'll need a working Erlang/Elixir environment on your host machine.
*   **Install AtomVM:** Follow AtomVM's build instructions, which involve cloning the repository and compiling.
*   **Flash:** Flash the AtomVM firmware to your ESP32/ESP8266.

### 4.2. Basic Display "Hello World" (Elixir)
*   Display drivers for AtomVM/Elixir on embedded platforms might be less mature or require more manual setup.
*   **Steps:**
    1.  Find or implement an AtomVM driver for ST7789/ST7735.
    2.  Write a simple Elixir application to initialize the display and draw text.
    3.  Compile and deploy your Elixir code to the ESP.

### 4.3. Wi-Fi Connectivity (Elixir)
*   AtomVM provides network capabilities.
*   **Steps:**
    1.  Implement Wi-Fi connection logic using AtomVM's network primitives.
    2.  Test HTTP requests.

---

## Decision Point:
After performing these initial setups and "Hello World" tests for each framework, you should be able to make an informed decision on which one to proceed with for the main project development. Consider:
*   Ease of setup and development for you.
*   Availability of display drivers and network libraries.
*   Community support.
*   Performance characteristics.

## Next Steps:
*   Based on your chosen framework, the subsequent tasks (`5_display_hello_world.md`, `6_wifi_connectivity.md`, etc.) will be tailored to that specific environment. You will then proceed with those specific instructions.
