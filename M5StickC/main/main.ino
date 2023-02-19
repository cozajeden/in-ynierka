#include <WiFi.h>
#include <EEPROM.h>
#include <M5StickCPlus.h>

#define pro_cpu 0
#define app_cpu 1
// pierwszy i drugi bit pierwszego bajtu 10 -un auth 01 - auth
// niezautoryzowany drugi = 0xFF
// zautoryzowany drugi bajt id
// odpowiedź pierwszy bajt:
//  statusOk = 1 
//  statusIdDoesntExist = 2
//  statusIdIsAlreadyConnected = 4
// odpowiedź drugi jeżeli prosiłem ID:
//  ID
//  jeżeli nie olewamy drugi bajt
const char* SSID = "HUAWEI-B525-2B90";
const char* PASSWORD =  "4MDH8DA9F6T";
 
const uint16_t PORT = 2137;
const char * HOST = "192.168.8.112";

const char* ntpServer = "time.google.com";
String timezone = "CET-1CEST,M3.5.0,M10.5.0/3";

WiFiClient client;

#define MAX_MSG_LEN 1436
byte msg[MAX_MSG_LEN];
int temp = -1;
const byte imu_header[] = {0x01, 0x01};

struct TimeUnionStruct {
  byte sec[4];
  byte usec[4];
};
union TimeUnion {
  struct TimeUnionStruct bytes;
  struct timeval tv_now;
};

void send_header() {
  Serial.println("Starting collectiong data from MPU6886");
  int data_cnt = 0;
  float* data_ptr = 0;
  TimeUnion timeUnion;
  memcpy(msg, imu_header, 2);
  msg[10] = 0xF0;
  msg[11] = 0x00;
  msg[12] = 0x00;
  msg[13] = 0x00;
  while(client.connected()){
    gettimeofday(&timeUnion.tv_now, NULL);
    timeUnion.tv_now.tv_sec += 7200;
    memcpy(msg + 2, timeUnion.bytes.sec, 8); // Accesing 8 bytes will write sec as well as usec
    data_cnt = 15;
    msg[14] = 0;
    while (data_cnt + 24 < MAX_MSG_LEN) {
      data_ptr = (float*)(msg + data_cnt);
      M5.Imu.getGyroData(data_ptr, data_ptr + 1, data_ptr + 2);
      M5.Imu.getAccelData(data_ptr + 3, data_ptr + 4, data_ptr + 5);
      data_cnt += 24;
      msg[14]++;
    }
    
//    Serial.print("sending data for MPU6886: ");
//    for (int i=0; i<data_cnt; i++) {
//      Serial.print(msg[i], HEX);
//    }
//    Serial.println('h');
    client.write(msg, MAX_MSG_LEN);
//    sleep(1000);
  }
  
  
}

static const int led = 10;
static const uint8_t queue_len = 5;

static SemaphoreHandle_t bin_sem;
static QueueHandle_t queue;

#define EEPROM_SIZE 1
const static uint8_t ID_address = 0;
byte ID = 0;
const byte ask_ID[] = {0x80, 0xFF};
byte send_ID[] = {0x40, 0x00};

int rec_msg() {
  temp = -1;
  while(temp == -1){
    // Wait for data
    temp = client.read();  
  }
  
  int i = 0;
  while(temp != -1){
    // Receive data
    msg[i] = temp;
    i++;
    temp = client.read();
  }
  return i;
}

void hand_shake() {
  bool Error = true;
  int len = 0;

  while (Error) {
    if (ID == 0xFF) {
      Serial.println("Asking for ID");
      client.write(ask_ID, 2);
      len = rec_msg();
      
      if (msg[0] == 1) {
        Serial.print("GOT statusOk with ID: ");
        Serial.print(msg[1], HEX);
        Serial.println('h');
        EEPROM.write(ID_address, msg[1]);
        EEPROM.commit();
        Serial.println("ID stored in Flash memory");
        Error = false;
      } else {
        Serial.print("Wrong response: ");
        for (int i=0; i < len; i++){
          Serial.print(msg[i], HEX);
        }
        Serial.println();
        sleep(1000);
        Serial.println("Trying again");
        Error = true;
      }
    } else {
      Serial.print("Sending ID: ");
      Serial.print(ID, HEX);
      Serial.println('h');
      send_ID[1] = ID;
      client.write(send_ID, 2);
      
      len = rec_msg();
  
      if (msg[0] == 1) {
        Serial.println("Received: statusOk");
        Error = false;
      } else if (msg[0] == 2) {
        Serial.println("Received: statusIdDoesntExist");
        Serial.println("Asking for new ID");
        EEPROM.write(ID_address, 0xFF);
        EEPROM.commit();
        Error = true;
      } else if (msg[0] == 4) {
        Serial.println("Received: statusIdIsAlreadyConnected");
        Serial.println("Trying again");
        Serial.println("Asking for new ID");
        EEPROM.write(ID_address, 0xFF);
        EEPROM.commit();
        Error = true;
      } else {
        Serial.print("Wrong response: ");
        for (int i=0; i < len; i++){
          Serial.print(msg[i], HEX);
        }
        Serial.println();
        Serial.println("Trying again");
        Error = true;
      }
    }
  }
}

String charArrayToString(char arrChar[], int tam) {
  String s = "";
  for (int i = 0; i < tam; i++) {
    s = s + arrChar[i];
  }
  return s;
}

void sleep(int ms) {
  vTaskDelay(ms/portTICK_PERIOD_MS);
}

bool connectClient() {
  if (!client.connect(HOST, PORT)) {
    Serial.println("Connection to host failed");
    sleep(1000);
    return false;
  }
  Serial.println("Connection to host successful");
  return true;
}

void connectWiFi() {
  Serial.println("Connecting to Wi-Fi");
  WiFi.begin(SSID, PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    sleep(500);
    Serial.print(".");
  }
  Serial.println("");
  
  Serial.print("WiFi connected with IP: ");
  Serial.println(WiFi.localIP());
}

void ntp_time() {
  configTime(3600, 3600, ntpServer);
  Serial.print("Setting up time zone to: ");
  Serial.println(timezone);
  setenv("TZ", timezone.c_str(),1);
  tzset();
  Serial.print("Trying to get time info from: ");
  Serial.println(ntpServer);
  Serial.println("It could take a while.");
}

void wait_for_time_info() {
  struct tm time;
  while (!getLocalTime(&time)) {
    Serial.println("Could not obtain time info");
  }
  Serial.println("Got current time");
}

void print_current_time_until_connected_to_the_server() {
  struct timeval tv_now;
  long int tt = tv_now.tv_usec;
  do {
    gettimeofday(&tv_now, NULL);
    tv_now.tv_sec += 7200;
    int year = tv_now.tv_sec/31556926;
    int day = (tv_now.tv_sec%31556926)/86400; 
    int hour = (tv_now.tv_sec%86400)/3600;
    int minu = (tv_now.tv_sec%3600)/60;
    int sec = tv_now.tv_sec%60;
    Serial.print("Years from 1970 to now: ");
    Serial.println(year);
    Serial.print("       day of the year: ");
    Serial.println(day);
    Serial.print("                 hours: ");
    Serial.println(hour);
    Serial.print("               minutes: ");
    Serial.println(minu);
    Serial.print("               seconds: ");
    Serial.println(sec);
  } while(!connectClient());
}

void setup() {
  Serial.begin(115200);
  sleep(7000);
  Serial.println("---START---");
  EEPROM.begin(EEPROM_SIZE);
  ID = EEPROM.read(ID_address);
//  M5.begin(); 
  M5.Imu.Init();
  queue = xQueueCreate(queue_len, sizeof(NULL));
  connectWiFi();
  ntp_time();
  wait_for_time_info();
  print_current_time_until_connected_to_the_server();
  hand_shake();
}

void loop() {
  send_header();
  if (!client.connected()) {
    print_current_time_until_connected_to_the_server();
    hand_shake();
  }
}
