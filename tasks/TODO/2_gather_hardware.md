# Task: Gather Required Hardware Components

## Description
This task involves physically gathering all the necessary hardware components for the project.

## Components List:
1.  **Microcontroller:**
    *   ESP32 (Recommended for more processing power and memory, easier display driving) OR
    *   ESP8266 (More limited resources, but might suffice for simple display and API calls).
2.  **Display:**
    *   160x128 ST7789 or ST7735 based display. (User provided)
3.  **Connectivity & Prototyping:**
    *   Breadboard
    *   Jumper wires (male-to-male, male-to-female, female-to-female as needed)
    *   USB-to-Serial converter (if using bare ESP modules without integrated USB-UART)
    *   Micro-USB cable (for powering ESP modules and flashing firmware)
4.  **Optional (for final product phase but good to have handy):**
    *   Solderable breadboard or perfboard
    *   Soldering iron and solder
    *   Multimeter (for debugging connections)

## Considerations:
*   Ensure the display's voltage requirements (3.3V or 5V) match the ESP's output. Most ESP modules operate at 3.3V logic. If your display is 5V, you'll need level shifters for the data lines. Many modern ST7789/ST7735 modules have built-in level shifting or are 3.3V compatible. Verify your specific display module.
*   Check the display's interface (SPI is common for these displays).
*   For the ESP32, specific pins might be better for SPI communication (VSPI/HSPI).
*   Ensure you have a reliable power supply for your breadboard and ESP.

## Next Steps:
1.  Verify the exact model of your ST7789/ST7735 display to confirm voltage compatibility and pinout.
2.  Ensure all listed components are available and in working order.
