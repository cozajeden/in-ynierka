#include "mike.hpp"

uint32_t Mike::sample_rate = 44100;
SemaphoreHandle_t mutexMike = xSemaphoreCreateBinary();

Mike::Mike() {
}

bool Mike::set_parameters(byte* params, int* length) {
  if (*length == 1) {
     if (params[0] == SampleFreq::Hz_22050) {
       sample_rate = 22050;
       return true;
     } else if (params[0] == SampleFreq::Hz_44100) {
       sample_rate = 44100;
       return true;
     } else {
       return false;
     }
  } else {
    return false;
  }
}

void Mike::init() {
    current_msg = MIKE_QUEUE_LEN;

    union {
        float var;
        byte bytes[4];
    } freq_u;
    freq_u.var = sample_rate;

    for (int i; i < MIKE_QUEUE_LEN; i++) {
        msg[i] = DataUnion();
        msg[i].data.header[0] = Commands::SEND_DATA;
        msg[i].data.header[1] = Devices::MIKE;
        memcpy(msg[i].data.frequency, freq_u.bytes, 4);
        msg_mutex[i] = xSemaphoreCreateBinary();
        xSemaphoreGive(msg_mutex[i]);
    }

    i2s_config_t cfg_i2s = {
        .mode = (i2s_mode_t)(I2S_MODE_MASTER | I2S_MODE_RX | I2S_MODE_PDM),
        .sample_rate = sample_rate,
        .bits_per_sample =
            I2S_BITS_PER_SAMPLE_16BIT,  // is fixed at 12bit, stereo, MSB
        .channel_format       = I2S_CHANNEL_FMT_ALL_RIGHT,
        #if ESP_IDF_VERSION > ESP_IDF_VERSION_VAL(4, 1, 0)
            .communication_format = I2S_COMM_FORMAT_STAND_I2S,
        #else
            .communication_format = I2S_COMM_FORMAT_I2S,
        #endif
        .intr_alloc_flags     = ESP_INTR_FLAG_LEVEL1,
        .dma_buf_count        = 2,
        .dma_buf_len          = 128,
    };
    


    i2s_pin_config_t cfg_pin;

#if (ESP_IDF_VERSION > ESP_IDF_VERSION_VAL(4, 3, 0))
    cfg_pin.mck_io_num = I2S_PIN_NO_CHANGE;
#endif

    cfg_pin.bck_io_num   = I2S_PIN_NO_CHANGE;
    cfg_pin.ws_io_num    = MIKE_PIN_CLK;
    cfg_pin.data_out_num = I2S_PIN_NO_CHANGE;
    cfg_pin.data_in_num  = MIKE_PIN_DATA;

    i2s_driver_install(I2S_NUM_0, &cfg_i2s, 0, NULL);
    i2s_set_pin(I2S_NUM_0, &cfg_pin);
    i2s_set_clk(I2S_NUM_0, sample_rate, I2S_BITS_PER_SAMPLE_16BIT, I2S_CHANNEL_MONO);
}

DataUnion* Mike::get_msg() {
    current_msg = (current_msg + 1) % MIKE_QUEUE_LEN;
    xSemaphoreTake(*get_current_semaphore(), portMAX_DELAY);
    return (msg + current_msg);
}

SemaphoreHandle_t* Mike::get_current_semaphore() {
    return (msg_mutex + current_msg);
}

void Mike::get_data(byte* data, byte* length) {
    size_t bytes_read = 0;
    i2s_read(I2S_NUM_0, (char *)data, MIKE_READ_LEN, &bytes_read, portMAX_DELAY);
    // 2 bajts mean 1 sample
    bytes_read = bytes_read / 2;
    length[0] = bytes_read & 0xFF;
    length[1] = (bytes_read >> 8) & 0xFF;
}

void mike_task(void *connection) {
    Connection * conn = (Connection*)connection;
    Mike mike = Mike();
    mike.init();
    DataUnion * msg = mike.get_msg();
    byte *data_ptr = (msg->msg.bytes + HEADER_LEN);
    TimeUnion timeUnion;
    Msg mm = Msg();
    float sample_interal = 1000000.0 / (float)mike.sample_rate;
    float time_passed = 0;
    gettimeofday(&timeUnion.tv_now, NULL);
    timeUnion.tv_now.tv_sec += 7200;
    while (uxSemaphoreGetCount(mutexMike) > 0) {
        memcpy(msg->data.time, timeUnion.bytes.sec, 8);
        mike.get_data(data_ptr, msg->data.count);
        time_passed = (sample_interal * (float)(msg->data.count[0] + (msg->data.count[1] << 8)));
        if (time_passed + timeUnion.tv_now.tv_usec > 1000000) {
            timeUnion.tv_now.tv_sec += 1;
            timeUnion.tv_now.tv_usec = time_passed + timeUnion.tv_now.tv_usec - 1000000;
        } else {
            timeUnion.tv_now.tv_usec += time_passed;
        }
        // Send data
        mm.data = msg->msg.bytes;
        mm.len = HEADER_LEN + ((msg->data.count[0] | (msg->data.count[1] << 8)) * 2);
        mm.mutex = mike.get_current_semaphore();
        conn->add_to_send_queue(&mm);
        // Get pointer to next free message
        msg = mike.get_msg();
        data_ptr = (msg->msg.bytes + HEADER_LEN);
    }
    i2s_driver_uninstall(I2S_NUM_0);
    sleep(50);
    xSemaphoreGive(mutexMike);
    while(true){}
}

bool start_Mike_task(Connection* connection) {
    xSemaphoreGive(mutexMike);
    xTaskCreatePinnedToCore(mike_task, "Mike", 8192, connection, configMAX_PRIORITIES -1, &taskMike, pro_cpu);
    return true;
}

bool stop_Mike_task() { 
    xSemaphoreTake(mutexMike, portMAX_DELAY); 
    xSemaphoreTake(mutexMike, portMAX_DELAY);
    if (taskMike != NULL) vTaskDelete(taskMike);
    return true;
}