#include "ESP8266HTTPClient.h"
#include "ESP8266WiFi.h"
#include "WiFiClient.h"
#include <Wire.h>
#include <Adafruit_SH110X.h>
#include <ArduinoJson.h>
#include "credentials.h"

HTTPClient http; 
WiFiClient wifiClient;

String playbackPath = String(SERVER_ORIGIN) + "/playback";

#define i2c_Address 0x3c

#define SCREEN_WIDTH 128
#define SCREEN_HEIGHT 64
#define OLED_RESET -1

Adafruit_SH1106G display = Adafruit_SH1106G(SCREEN_WIDTH, SCREEN_HEIGHT, &Wire, OLED_RESET);

unsigned long lastTime = 0;
unsigned long timerDelay = 1000;

void setup() {
  Serial.begin(115200);
  WiFi.begin(SSID, SSID_PASSWORD);
 
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

    if (doc["status"] == "inactive") {
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
  display.println("Spotify nic nie gra");
}

void handleActivePlayer(StaticJsonDocument<1000> doc) {
  display.println("Spotify gra");

  String artist = doc["artists"];
  String title = doc["name"];
  int progress = doc["progressMs"];
  int duration = doc["durationMs"];

  display.println(artist);
  display.println(title);
  display.println(progress);
  display.println(duration);
}
