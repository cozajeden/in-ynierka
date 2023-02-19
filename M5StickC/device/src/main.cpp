#include "globals.hpp"
#include "link.hpp"
#include "mpu.hpp"
#include "mike.hpp"

Connection conn = Connection();

SemaphoreHandle_t command_response_mutex;

static const int led = 10;
int led_delay = 2000;
static const uint8_t queue_len = 5;
byte streaming_device = Devices::NONE;

void Write1Byte(uint8_t Addr, uint8_t Data) {
    Wire1.beginTransmission(0x34);
    Wire1.write(Addr);
    Wire1.write(Data);
    Wire1.endTransmission();
}

uint8_t Read8bit(uint8_t Addr) {
    Wire1.beginTransmission(0x34);
    Wire1.write(Addr);
    Wire1.endTransmission();
    Wire1.requestFrom(0x34, 1);
    return Wire1.read();
}

void power_up(){
    Wire1.begin(21, 22);
    Wire1.setClock(400000);

    // Set ADC to All Enable
    Write1Byte(0x82, 0xff);

    // Bat charge voltage to 4.2, Current 100MA
    Write1Byte(0x33, 0xc0);

    // Enable Bat,ACIN,VBUS,APS adc
    Write1Byte(0x82, 0xff);

    // Enable Ext, LDO2, LDO3, DCDC1
    Write1Byte(0x12, Read8bit(0x12));

    // 128ms power on, 4s power off
    Write1Byte(0x36, 0x0C);

    // Set GPIO0 to LDO
    Write1Byte(0x90, 0x02);

    // Disable vbus hold limit
    Write1Byte(0x30, 0x80);

    // Set temperature protection
    Write1Byte(0x39, 0xfc);

    // Enable bat detection
    Write1Byte(0x32, 0x46);
}

bool device_on(byte device, byte *response)
{
  D(Serial.println("Device on"););
  response[1] = Responses::COMMAND_OK;
  switch(device)
  {
    case Devices::MPU6886:
    D(Serial.println("MPU on"););
      return start_MPU6886_task(&conn);
    case Devices::MIKE:
    D(Serial.println("Device on: Mike"););
      return start_Mike_task(&conn);
    // Here add other devices
    // case Devices::Example:
    //  D(Serial.println("Device on: Example"););
    //  return start_Example_task(&conn);
    default:
      response[1] = Responses::COMMAND_UNKNOWN;
      return false;
  }
}

bool device_off(byte device, byte *response)
{
  D(Serial.println("Device off"););
  response[1] = Responses::COMMAND_OK;
  switch(device)
  {
    case Devices::MPU6886:
    D(Serial.println("MPU6886 off"););
      return stop_MPU6886_task();
    case Devices::MIKE:
    D(Serial.println("MIKE off"););
      return stop_Mike_task();
    // Here add other devices
    // case Devices::Example:
    //  D(Serial.println("Example off"););
    //  return stop_Example_task();
    default:
      response[1] = Responses::COMMAND_UNKNOWN;
      return false;
  }
}

bool set_device_parameters(byte device, byte *params, int *param_len)
{
  D(Serial.println("Set device parameters"););
  switch (device) {
    case Devices::MPU6886:
      D(Serial.println("Setting MPU parameters");)
      return MPU::set_parameters(params, param_len);
    case Devices::MIKE:
      D(Serial.println("Setting MIKE parameters");)
      return Mike::set_parameters(params, param_len);
    // Here add other devices
    // case Devices::Example:
    //  D(Serial.println("Setting Example parameters");)
    //  return Example::set_parameters(params, param_len);
    default:
      D(Serial.println("Command: 0x?? - (Unknown)");)
      return false;
  }
}

void send_task(void *parameter)
{
  D(Serial.println("Send task"););
  // Task responsible for sending data to the server
  Connection *conn = (Connection *)parameter;
  while(true){
      conn->send_from_queue();
  }
}

bool device_parameters(byte *command, int *command_len, byte *response)
{
  D(Serial.println("Device parameters"););
  *command_len -= 2;
  if (set_device_parameters(command[1], &command[2], command_len)) {
    response[1] = Responses::COMMAND_OK;
    return true;
  } else {
    response[1] = Responses::COMMAND_UNKNOWN;
    D(Serial.println("Wrong parameters");)
    return false;
  }
}

bool on_off(byte *command, int *command_len, byte *response)
{
  D(Serial.println("on_off");)
  // Run command
  switch(command[2]) {
    case Commands::OFF:
      D(Serial.println("Turning off");)
      if (streaming_device == Devices::NONE) {
        response[1] = Responses::COMMAND_NO_STREAMING;
      } else {
        if(device_off(command[1], response)) {
          streaming_device = Devices::NONE;
          return true;
        }
      }
      break;
    case Commands::ON:
      D(Serial.println("Turning on");)
      if (streaming_device != Devices::NONE) {
        response[1] = Responses::COMMAND_ALREADY_STREAMING;
      } else {
        if( device_on(command[1], response)) {
          streaming_device = command[1];
          return true;
        }
      }
      break;
  }
  return false;
}

bool run_command(byte *command, int *command_len, byte *response)
{
  D(Serial.println("run_command");)
  // Check minimum length of command
  if(*command_len < 3) {
    D(Serial.println("Command too short");)
    if ((*command_len) > 0) {
      response[0] = command[0];
    } 
    response[1] = Responses::COMMAND_UNKNOWN;
    sleep(100);
    conn.hand_shake();
    return false;
  } 
  response[0] = command[0];
  // Check which command is being sent
  switch(command[0]) {
    case Commands::SEND_PARAMETERS:
      return device_parameters(command, command_len, response);
    case Commands::ON_OFF:
      return on_off(command, command_len, response);
    default:
      D(Serial.println("Command: 0x?? - (Unknown)");)
      response[1] = Responses::COMMAND_UNKNOWN;
      return false;
  }
}

void led_task(void *parameter)
{
  while (true) {
    digitalWrite(led, LOW);
    sleep(10);
    digitalWrite(led, HIGH);
    sleep(led_delay);
    if (streaming_device != Devices::NONE) {
    digitalWrite(led, LOW);
    sleep(10);
    digitalWrite(led, HIGH);
    }
    sleep(50);
  }
}

void main_task(void *parameter)
{
  Connection *conn = (Connection *)parameter;
  xTaskCreatePinnedToCore(send_task, "send_task", 4096, conn, 2, NULL, app_cpu);
  xTaskCreatePinnedToCore(led_task, "led_task", 1024, &conn, 1, NULL, app_cpu);
  
  Msg msg = Msg();
  byte data[2] = {0x00, 0x00};
  msg.data = data;
  msg.len = 2;
  msg.mutex = &command_response_mutex;
  int command_max_len = 4;
  int received_bytes = 0;
  byte command[command_max_len];
  bool result = false;
  xSemaphoreGive(command_response_mutex);

  while (true) {
    conn->rec_msg(command, &received_bytes);
    run_command(command, &received_bytes, msg.data);
    // If you want to send a response to the command, uncomment the following lines
      // xSemaphoreTake(command_response_mutex, portMAX_DELAY);
      // D(Serial.println("Sending response");)
      // conn->add_to_send_queue(&msg);
  }
}

void setup()
{
  D(
    Serial.begin(115200);
    sleep(500);
    Serial.println("---START---");
  );
  EEPROM.begin(EEPROM_SIZE);
  pinMode(led, OUTPUT);
  command_response_mutex = xSemaphoreCreateBinary();
  power_up();
  conn.init();
  xTaskCreatePinnedToCore(main_task, "main_task", 4096, &conn, 1, NULL, app_cpu);
}

void loop(){}