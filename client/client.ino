#include "ESP8266HTTPClient.h"
#include "ESP8266WiFi.h"
#include "WiFiClient.h"
#include <Wire.h>
#include <Adafruit_SH110X.h>
#include <ArduinoJson.h>
#include "credentials.h"

HTTPClient http; 
WiFiClientSecure wifiClient;

String playbackPath = String(SERVER_ORIGIN) + "/playback";

#define i2c_Address 0x3c

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_RESET -1

Adafruit_SH1106G display = Adafruit_SH1106G(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);

unsigned long lastTime = 0;
unsigned long timerDelay = 1000;

int SpotifyLogoSize = 28;
const unsigned char SpotifyLogo [] PROGMEM = {
  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
  0x00, 0x0f, 0xf0, 0x00, 0x00, 0x3f, 0xfc, 0x00, 0x00, 0xff, 0xff, 0x00, 0x01, 0xff, 0xff, 0x80,
  0x03, 0xff, 0xff, 0xc0, 0x03, 0xff, 0xff, 0xc0, 0x07, 0x00, 0x07, 0xe0, 0x06, 0x00, 0x01, 0xe0,
  0x0e, 0x00, 0x00, 0x70, 0x0f, 0xff, 0xf0, 0x30, 0x0f, 0xc0, 0x3e, 0x70, 0x0f, 0x00, 0x07, 0xf0,
  0x0f, 0x00, 0x01, 0xf0, 0x0f, 0xff, 0xe0, 0xf0, 0x0f, 0xc0, 0x79, 0xf0, 0x0f, 0x00, 0x0f, 0xf0,
  0x07, 0x86, 0x03, 0xe0, 0x07, 0xff, 0xe3, 0xe0, 0x03, 0xff, 0xff, 0xc0, 0x03, 0xff, 0xff, 0xc0,
  0x01, 0xff, 0xff, 0x80, 0x00, 0xff, 0xff, 0x00, 0x00, 0x3f, 0xfc, 0x00, 0x00, 0x0f, 0xf0, 0x00,
  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
};

void setup() {
  Serial.begin(115200);
  WiFi.begin(SSID, SSID_PASSWORD);
  wifiClient.setInsecure();
 
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.println("Waiting to connect…");
  }
 
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  display.begin(i2c_Address, true);
  display.setTextSize(1);
  display.setTextColor(SH110X_WHITE);
}
 
void loop() {
  bool connected = WiFi.status() == WL_CONNECTED;
  bool timeElapsed = (millis() - lastTime) > timerDelay;

  if (connected && timeElapsed) {
    http.begin(wifiClient, playbackPath);
    http.addHeader("boardToken", String(BOARD_TOKEN));

    int httpCode = http.GET();
    String resp = http.getString();

    StaticJsonDocument<1000> doc;
    DeserializationError error = deserializeJson(doc, resp);
    if (error) {
      Serial.print(F("deserializeJson() failed: "));
      Serial.println(error.f_str());
    }

    display.clearDisplay();
    display.setCursor(0, 0);

    if (doc["state"] == "inactive") {
      handleInactivePlayer();
    } else {
      handleActivePlayer(doc);
    }

    display.display();

    http.end();

    lastTime = millis();
  }
}

void handleInactivePlayer() {
  int x = 15;
  int y = SCREEN_HEIGHT / 2 - SpotifyLogoSize / 2;
  
  display.drawBitmap(x, y, SpotifyLogo, SpotifyLogoSize, SpotifyLogoSize, SH110X_WHITE);

  display.setCursor(x + SpotifyLogoSize + 10, y + 8);
  display.println("Spotify");
  display.setCursor(x + SpotifyLogoSize + 10, y + 16);
  display.println("turned off");
}

void handleActivePlayer(StaticJsonDocument<1000> doc) {
  display.println();
  display.setTextColor(SH110X_BLACK, SH110X_WHITE);
  if (doc["state"] == "playing") {
    display.println("Spotify plays");
  } else {
    display.println("Spotify paused");
  }

  display.println();

  display.drawLine(0, 20, 86, 20, SH110X_WHITE);

  display.setTextColor(SH110X_WHITE);

  String artist = doc["artists"];
  String title = doc["name"];
  float progress = doc["progressMs"];
  float duration = doc["durationMs"];

  display.println(title);

  display.setTextSize(1);
  display.println("by: " + artist);

  int width = progress / duration * SCREEN_WIDTH;
  int barHeight = 4;

  display.drawRect(0, SCREEN_HEIGHT - barHeight, width, barHeight, SH110X_WHITE);
  display.fillRect(0, SCREEN_HEIGHT - barHeight, width, barHeight, SH110X_WHITE);
  display.drawBitmap(SCREEN_WIDTH - SpotifyLogoSize, -4, SpotifyLogo, SpotifyLogoSize, SpotifyLogoSize, SH110X_WHITE);
}
