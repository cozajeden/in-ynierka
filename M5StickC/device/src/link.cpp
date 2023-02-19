#include "link.hpp"

Connection::Connection() {
  ID = 0;
  send_ID[0] = Commands::SEND_ID;
  send_ID[1] = 0x00;
  ask_ID[0] = Commands::ASK_ID;
  ask_ID[1] = 0xFF;
  recv_task = NULL;
}

void Connection::init() {
  ID = EEPROM.read(ID_EEPROM_ADDRESS);
  connectWiFi();
  ntp_time();
  wait_for_time_info();
  block_until_client_is_connected();
  hand_shake();
  send_queue = xQueueCreate(5, sizeof(struct Msg));
  recv_mutex = xSemaphoreCreateBinary();
  xSemaphoreGive(recv_mutex);
}

void Connection::add_to_send_queue(Msg* data) {
  xQueueSend(send_queue, data, portMAX_DELAY);
}

void Connection::send_from_queue() {
  Msg msg = Msg();
  // D(Serial.println(msg.len););
  xQueueReceive(send_queue, &msg, portMAX_DELAY);
  if (!client.connected()) {
    block_until_client_is_connected();
    hand_shake();
  }
  client.write(msg.data, msg.len);
  xSemaphoreGive(*msg.mutex);
}

bool Connection::block_until_client_is_connected() {
  if (!client.connected()) {
  while (!connectClient()) {};
  return false;
  }
  return true;
}

void Connection::rec_msg(byte* msg, int* len) {
  size_t bytes_to_read = 0;
  *len = 0;
  byte* data_ptr = msg;

  while (bytes_to_read == 0) {
    if (!block_until_client_is_connected()) {
      data_ptr[0] = 0x00;
      return;
    }
    bytes_to_read = client.available();
  }

  while (bytes_to_read != 0) {
    client.read((data_ptr + *len), bytes_to_read);
    *len += bytes_to_read;
    bytes_to_read = client.available();
  }
}

bool Connection::connectClient() {
  if (!client.connect(HOST, PORT)) {
    D(Serial.println("Connection to host failed");)
    sleep(1000);
    return false;
  }
  D(Serial.println("Connection to host successful"););
  return true;
}

void Connection::hand_shake() {
  bool Error = true;
  int len = 0;


  while (Error) {
    block_until_client_is_connected();
    ID = EEPROM.read(ID_EEPROM_ADDRESS);
    if (ID == 0xFF) {
      D(Serial.println("Asking for ID"););
      client.write(ask_ID, 2);
      if (!client.connected()) continue;
      rec_msg(msg, &len);

      if (msg[0] == Responses::HANDSHAKE_OK) {
        EEPROM.write(ID_EEPROM_ADDRESS, msg[1]);
        EEPROM.commit();
        D(
          Serial.print(msg[1], HEX);
          Serial.println('h');
          Serial.print("GOT statusOk with ID: ");
          Serial.println("ID stored in Flash memory");
        );
        Error = false;
      } else {
        D(
          Serial.print("Wrong response: ");
          for (int i = 0; i < len; i++) {
            Serial.print(msg[i], HEX);
          }
          Serial.println();
          Serial.println("Trying again");
        );
        Error = true;
        sleep(1000);
      }
    } else {
      D(
        Serial.print("Sending ID: ");
        Serial.print(ID, HEX);
        Serial.println('h');
      );
      send_ID[1] = ID;
      client.write(send_ID, 2);

      if (!client.connected()) continue;
      rec_msg(msg, &len);
      switch (msg[0])
        {
          case Responses::HANDSHAKE_OK:
            D(Serial.println("Received: statusOk"););
            Error = false;
            break;
          case Responses::HANDSHAKE_ID_DOESNT_EXIST:
            D(Serial.println("Received: IdDoesntExist"););
            Error = true;
            break;
          case Responses::HANDSHAKE_ID_IS_ALREADY_CONNECTED:
            D(Serial.println("Received: IdIsAlreadyConnected"););
            Error = true;
            break;
          case Responses::HANDSHAKE_NO_ID_AVAILABLE:
            D(Serial.println("Received: NoIdAvailable"););
            Error = true;
            break;
          case Responses::HANDSHAKE_HANDSHAKE_NOT_VALID:
            D(Serial.println("Received: HandshakeNotValid"););
            Error = true;
            break;
          case Responses::HANDSHAKE_COMMAND_NOT_RECOGNIZED:
            D(Serial.println("Received: CommandNotRecognized"););
            Error = true;
            break;
          case Responses::HANDSHAKE_UNEXPECTED_PACKET_HEADER:
            D(Serial.println("Received: UnexpectedPacketHeader"););
            Error = true;
            break;
          case Responses::HANDSHAKE_UNEXPECTED_PACKET_BODY:
            D(Serial.println("Received: UnexpectedPacketBody"););
            Error = true;
            break;
          
          default:
            D(
              Serial.print("Wrong response: ");
              for (int i = 0; i < len; i++) {
                Serial.print(msg[i], HEX);
              }
            );
            Error = true;
            continue;
        } 
      if (Error) {
        // If any error - set ID to unset and try again 
        EEPROM.write(ID_EEPROM_ADDRESS, 0xFF);
        EEPROM.commit();
      }
    }
  }
}

void Connection::connectWiFi() {
  D(Serial.println("Connecting to Wi-Fi"););
  WiFi.begin(SSID, PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    sleep(500);
    D(Serial.print("."););
  }
  D(
    Serial.println("");
    Serial.print("WiFi connected with IP: ");
    Serial.println(WiFi.localIP());
  );
}

void Connection::ntp_time() {
  configTime(3600, 3600, ntpServer);
  D(
    Serial.print("Trying to get time info from: ");
    Serial.println(ntpServer);
    Serial.println("It could take a while.");
  );
}

void Connection::wait_for_time_info() {
  struct tm time;
  while (!getLocalTime(&time)) {
    D(Serial.println("Could not obtain time info"););
  }
  D(Serial.println("Got current time"););
}
