#include "mpu.hpp"

SemaphoreHandle_t mutex6886 = xSemaphoreCreateBinary();

MPU::MPU(){}
byte MPU::settings[2] = {0x00, 0x00};
TwoWire MPU::i2c = TwoWire(0);

bool MPU::set_parameters(byte* params, int* length) {
    if (*length != 2) {
        return false;
    }
    if (params[0] >= MPUGyroRange::DPS_250 && params[0] <= MPUGyroRange::DPS_2000 && params[1] >= MPUAccelRange::G2 && params[1] <= MPUAccelRange::G16) {
        settings[0] = params[0];
        settings[1] = params[1];
        return true;
    } else {
        return false;
    }
}

void MPU::init() {
    init(settings[1], settings[0]);
}

void MPU::init(byte a_range, byte g_range) {
    D(Serial.println("Initializing MPU"););
    current_msg = MPU_QUEUE_LEN;
    for (int i; i < MPU_QUEUE_LEN; i++) {
        msg[i] = DataUnion();
        msg[i].data.header[0] = Commands::SEND_DATA;
        msg[i].data.header[1] = Devices::MPU6886;
        msg_mutex[i] = xSemaphoreCreateBinary();
        xSemaphoreGive(msg_mutex[i]);
    }
    D(Serial.println("Initializing MPU complete"););
    D(Serial.println("Initializing I2C"););
    i2c.begin(MPU_SDA_PIN, MPU_SCL_PIN);
    D(Serial.println("Initializing I2C complete"););
    reset();
    D(Serial.println("Resetting MPU complete"););
    set_accel_range(a_range);
    set_gyro_range(g_range);
    accel_range_f = get_accel_range();
    gyro_range_f = get_gyro_range();
    accel_multiplier = accel_range_f/32786.0;
    gyro_multiplier  = gyro_range_f/32786.0;
    D(Serial.println("Initializing MPU complete"););
    set_frequency();
    D(Serial.println("Measuring MPU frequency complete"););
}

void MPU::reset() {
    unsigned char regdata = 0x00;
    write_NBytes(MPU_ADDRESS, MPU_PWR_MGMT_1, 1, &regdata);
    sleep(10);
}

void MPU::write_NBytes(uint8_t driver_Address, uint8_t reg_Address, uint8_t length, uint8_t* data) {
    i2c.beginTransmission(driver_Address);
    i2c.write(reg_Address);
    i2c.write(data, length);
    i2c.endTransmission(true);
}

void MPU::read_NBytes(uint8_t driver_Address, uint8_t reg_Address, uint8_t length, uint8_t* data) {
    i2c.beginTransmission(driver_Address);
    i2c.write(reg_Address);
    i2c.endTransmission(false);
    i2c.requestFrom(driver_Address, length);
    uint8_t i = 0;
    while (i2c.available()) {
        data[i++] = i2c.read();
    }
}

void MPU::get_data(float* data_ptr) {
    read_NBytes(MPU_ADDRESS, MPU_ACCEL_XOUT_H, 14, data);
    for (int i = 0; i < 3; i++) {
        calculate_from_bytes(
            2*i,
            &accel_multiplier,
            &accel_range_f,
            (data_ptr + i)
        );
        calculate_from_bytes(
            2*i+8,
            &gyro_multiplier,
            &gyro_range_f,
            (data_ptr + 3 + i)
        );
    }
}

DataUnion* MPU::get_msg() {
    current_msg = (current_msg + 1) % MPU_QUEUE_LEN;
    xSemaphoreTake(*get_current_semaphore(), portMAX_DELAY);
    return (msg + current_msg);
}

SemaphoreHandle_t* MPU::get_current_semaphore() {
    return (msg_mutex + current_msg);
}

void MPU::set_frequency() {
    union {
        float var;
        byte bytes[4];
    } freq_u;
    

    freq_u.var = 582.0;
    
    for (int i=0; i < MPU_QUEUE_LEN; i++) {
        for (int j=0; j < 4; j++) {
            msg[i].msg.bytes[j+10] = freq_u.bytes[j];
            D(Serial.print(freq_u.bytes[j], HEX));
        }
        D(Serial.println(""););
    }
}

void MPU6886_task(void* connection) { 
    Connection * conn = (Connection*)connection;

    // Prepare the MPU6886
    MPU mpu = MPU();
    mpu.init();
    DataUnion * msg = mpu.get_msg();
    float *data_ptr = (float *)(msg->msg.bytes + HEADER_LEN);
    int max_len = (MAX_MSG_LEN - HEADER_LEN)/24;
    size_t len = 0;
    TimeUnion timeUnion;
    Msg mm = Msg();

    // Main measurement loop
    while (uxSemaphoreGetCount(mutex6886) > 0) {
        // Set the time of the first measurement
        len = 0;
        gettimeofday(&timeUnion.tv_now, NULL);
        timeUnion.tv_now.tv_sec += 7200;
        memcpy(msg->data.time, timeUnion.bytes.sec, 8);
        // Take measurements
        while (len < max_len) {
            mpu.get_data((data_ptr + len*6));
            len += 1;
        }
        // Send data
        mm.data = msg->msg.bytes;
        msg->data.count[0] = (len) & 0xFF;
        msg->data.count[1] = (len >> 8) & 0xFF;
        mm.len = HEADER_LEN + len*24;
        mm.mutex = mpu.get_current_semaphore();
        conn->add_to_send_queue(&mm);
        // Get pointer to next free message
        msg = mpu.get_msg();
        data_ptr = (float *)(msg->msg.bytes + HEADER_LEN);
    }
    D(Serial.println("Ending Wire"););
    mpu.i2c.end();
    D(Serial.println("Wire ended"););
    xSemaphoreGive(mutex6886);
    while(true){}
}

bool start_MPU6886_task(Connection *connection) {
    D(Serial.println("Starting MPU6886 task"););
    xSemaphoreGive(mutex6886);
    xTaskCreatePinnedToCore(MPU6886_task, "mpu_task", 8192, connection, configMAX_PRIORITIES -1, &task6886, pro_cpu);
    D(Serial.println("Starting MPU6886 task complete"););
    return true;
}

bool stop_MPU6886_task() {
    D(Serial.println("Stopping MPU6886 task"););
    D(Serial.println("Sending stop signal to MPU6886"););   
    xSemaphoreTake(mutex6886, portMAX_DELAY);
    D(Serial.println("Sent stop signal to MPU6886"););
    D(Serial.println("Waiting for MPU6886 task to end"););
    xSemaphoreTake(mutex6886, portMAX_DELAY);
    D(Serial.println("MPU6886 task ended"););
    if (task6886 != NULL) vTaskDelete(task6886);
    D(Serial.println("MPU6886 task stopped"););
    return true;
}

void MPU::set_accel_range(byte range) {
    unsigned char regdata = range;
    regdata = regdata << 3;
    write_NBytes(MPU_ADDRESS, MPU_ACCEL_CONFIG, 1, &regdata);
    sleep(50);
    accel_range = range;
}

void MPU::set_gyro_range(byte range) {
    unsigned char regdata = range;
    regdata = regdata << 3;
    write_NBytes(MPU_ADDRESS, MPU_GYRO_CONFIG, 1, &regdata);
    sleep(50);
    gyro_range = range;
}

float MPU::get_accel_range() {
    switch (accel_range) {
        case MPUAccelRange::G2:
            return 2.0;
        case MPUAccelRange::G4:
            return 4.0;
        case MPUAccelRange::G8:
            return 8.0;
        case MPUAccelRange::G16:
            return 16.0;
        default:
            return 2.0;
    }
}

float MPU::get_gyro_range() {
    switch (gyro_range) {
        case MPUGyroRange::DPS_250:
            return 250.0;
        case MPUGyroRange::DPS_500:
            return 500.0;
        case MPUGyroRange::DPS_1000:
            return 1000.0;
        case MPUGyroRange::DPS_2000:
            return 2000.0;
        default:
            return 250.0;
    }
}

void MPU::calculate_from_bytes(uint8_t start_index, float* scale, float* range, float* destination) {
    // Calculate the float value from two bytes obtained from the MPU6886
    if ((data[start_index] & 0x80) != 0x80) {
        *destination = (((data[start_index] & 0x7F) << 8) | data[start_index+1]) * *scale;
    } else {
        *destination =  -(~(((data[start_index] & 0x7F) << 8) | data[start_index+1])) * *scale - *range;
    }
}