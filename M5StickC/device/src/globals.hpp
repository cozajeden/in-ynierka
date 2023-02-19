#ifndef GLOBALS_HPP
#define GLOBALS_HPP

#include <Arduino.h>
#include <WiFi.h>
#include <EEPROM.h>
#include <M5StickCPlus.h>
#include <Wire.h>
#include <driver/i2s.h>

#define MAX_MSG_LEN 1436
#define HEADER_LEN 16
#define EEPROM_SIZE 1

#define pro_cpu 0
#define app_cpu 1

#define DEBUG 0

#if DEBUG == 1
#define D(x) x
#else 
#define D(x)
#endif
struct Devices {
  const static byte NONE {0x00};
  const static byte MPU6886 {0x01};
  const static byte MIKE {0x02};
};

struct Msg
{
    byte *data;
    int len;
    SemaphoreHandle_t* mutex;
};

struct TimeUnionStruct
{
  byte sec[4];
  byte usec[4];
};

union TimeUnion
{
  struct TimeUnionStruct bytes;
  struct timeval tv_now;
};

struct DataStruct
{
    byte header[2];
    byte time[8];
    byte frequency[4];
    byte count[2];
    byte data[MAX_MSG_LEN-HEADER_LEN];
};

struct MsgStruct
{
    byte bytes[MAX_MSG_LEN];
};

union DataUnion
{
    struct DataStruct data;
    struct MsgStruct msg;
};

String charArrayToString(char arrChar[], int tam);
void sleep(int ms);

#endif