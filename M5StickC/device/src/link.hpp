#pragma once
#ifndef __LINK__
#define __LINK__
#include "globals.hpp"

#define ID_EEPROM_ADDRESS 0
// #define SSID "Mi 10T"
// #define PASSWORD "grundfos"
// #define HOST "192.168.43.36"
#define SSID "HUAWEI-B525-2B90"
#define PASSWORD "4MDH8DA9F6T"
#define HOST "192.168.8.112"
#define PORT 7122

#if DEBUG == 1
#define ntpServer "time.google.com"
#else 
#define ntpServer HOST
#endif

struct Commands {
  const static byte SEND_DATA {0x03};
  const static byte ASK_ID {0x01};
  const static byte SEND_ID {0x02};
  const static byte SEND_PARAMETERS {0x04};
  const static byte ON_OFF {0x05};
  const static byte ON {0x01};
  const static byte OFF {0x02};
};

struct Responses {
  const static byte HANDSHAKE_OK {0x01};
  const static byte HANDSHAKE_ID_DOESNT_EXIST {0x02};
  const static byte HANDSHAKE_ID_IS_ALREADY_CONNECTED {0x03};
  const static byte HANDSHAKE_NO_ID_AVAILABLE {0x04};
  const static byte HANDSHAKE_HANDSHAKE_NOT_VALID {0x05};
  const static byte HANDSHAKE_COMMAND_NOT_RECOGNIZED {0x06};
  const static byte HANDSHAKE_UNEXPECTED_PACKET_HEADER {0x07};
  const static byte HANDSHAKE_UNEXPECTED_PACKET_BODY {0x08};
  const static byte COMMAND_UNKNOWN {0x0B};
  const static byte COMMAND_OK {0x01};
  const static byte COMMAND_ALREADY_STREAMING {0x09};
  const static byte COMMAND_NO_STREAMING {0x0A};
};

static TaskHandle_t recv_task;
void receiver_task(void* queue);

class Connection {
  
  private:
    int ID;
    byte send_ID[2];
    byte ask_ID[2];
    SemaphoreHandle_t recv_mutex;
    QueueHandle_t send_queue;
  
  public:
    Connection();
    WiFiClient client;
    byte msg[MAX_MSG_LEN];
    void rec_msg(byte* msg, int* len);
    void add_to_send_queue(Msg* data);
    void send_from_queue();
    bool connectClient();
    void hand_shake();
    void connectWiFi();
    void ntp_time();
    void wait_for_time_info();
    void init();
    bool block_until_client_is_connected();
    void start_receiver_task(QueueHandle_t* queue);
    void stop_receiver_task();
};
#endif
