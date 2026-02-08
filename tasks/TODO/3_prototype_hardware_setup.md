# Task: Prototype Phase - Hardware Setup on Breadboard

## Description
This task details the physical wiring of the chosen ESP module (ESP32 or ESP8266) to the ST7789/ST7735 display on a breadboard. Correct connections are crucial for the display to function.

## Prerequisites
*   Gathered hardware components as per `2_gather_hardware.md`.
*   Familiarity with your specific display module's pinout (ST7789/ST7735 displays can have slight variations in pin labels).
*   Multimeter for verifying connections.

## General Wiring Guidelines (SPI Interface)

Most ST7789/ST7735 displays communicate via SPI. Here's a common pin mapping. **Always refer to your specific display module's datasheet/pinout.**

| Display Pin (Common) | ESP32 Pin (Typical SPI VSPI) | ESP8266 Pin (Typical SPI HSPI) | Description                                       |
| :------------------- | :--------------------------- | :----------------------------- | :------------------------------------------------ |
| VCC                  | 3.3V                         | 3.3V                           | Power Supply (Connect to 3.3V on ESP)             |
| GND                  | GND                          | GND                            | Ground                                            |
| SCL (SCK)            | GPIO 18 (VSPI SCK)           | GPIO 14 (HSPI SCK)             | Serial Clock                                      |
| SDA (MOSI)           | GPIO 23 (VSPI MOSI)          | GPIO 13 (HSPI MOSI)            | Master Out Slave In (Data from ESP to Display)    |
| RES (RST)            | GPIO 4 (or any available GPIO) | GPIO 2 (or any available GPIO) | Reset Pin                                         |
| DC (A0, D/C)         | GPIO 2 (or any available GPIO) | GPIO 0 (or any available GPIO) | Data/Command Select (High for Data, Low for Command) |
| CS (CE)              | GPIO 5 (VSPI CS)             | GPIO 15 (HSPI CS)              | Chip Select                                       |
| BLK (LED)            | 3.3V (or GPIO with resistor) | 3.3V (or GPIO with resistor) | Backlight (usually powered by 3.3V directly, or via GPIO for brightness control) |

### Important Notes:
*   **Voltage Compatibility:** Double-check if your display requires 5V. If so, you MUST use a logic level shifter for all data lines (SCL, SDA, RES, DC, CS) to protect your 3.3V ESP module. Most modern ST7789/ST7735 boards are 3.3V compatible, but confirm.
*   **ESP32 SPI:** The ESP32 has two SPI interfaces (VSPI and HSPI). VSPI is generally preferred for external devices. The typical VSPI pins are SCK (GPIO18), MISO (GPIO19), MOSI (GPIO23), CS (GPIO5). Since the display is output-only, MISO isn't strictly needed.
*   **ESP8266 SPI:** The typical HSPI pins are SCK (GPIO14), MISO (GPIO12), MOSI (GPIO13), CS (GPIO15).
*   **Other GPIOs:** RES and DC pins can usually be assigned to any available GPIO pin on the ESP module.
*   **Backlight (BLK/LED):** Often, this can be connected directly to 3.3V for full brightness. If you need brightness control, connect it to a PWM-capable GPIO via a current-limiting resistor (check your display's datasheet for resistor value). For simplicity in prototyping, direct 3.3V is fine.

## Step-by-Step Wiring:
1.  **Mount:** Place the ESP module and the ST7789/ST7735 display on the breadboard. Ensure enough space for connections.
2.  **Power:** Connect the display's VCC to the ESP's 3.3V pin and GND to the ESP's GND.
3.  **SPI Bus:**
    *   Connect Display SCL to ESP's SPI Clock pin (e.g., ESP32: GPIO18, ESP8266: GPIO14).
    *   Connect Display SDA to ESP's SPI MOSI pin (e.g., ESP32: GPIO23, ESP8266: GPIO13).
4.  **Control Pins:**
    *   Connect Display RES to a chosen ESP GPIO pin (e.g., ESP32: GPIO4).
    *   Connect Display DC to a chosen ESP GPIO pin (e.g., ESP32: GPIO2).
    *   Connect Display CS to a chosen ESP SPI CS pin (e.g., ESP32: GPIO5, ESP8266: GPIO15).
5.  **Backlight:** Connect Display BLK/LED to 3.3V.

## Verification:
*   Visually inspect all connections to ensure they are firm and correct.
*   Double-check for any short circuits.
*   Do NOT power on yet until you have cross-referenced your display's specific pinout with these general guidelines.

## Next Steps:
*   Once wiring is complete and verified, proceed to `4_flash_micropython.md`.
