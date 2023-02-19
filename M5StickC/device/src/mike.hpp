#ifndef _MIKE_
#define _MIKE_
#include "globals.hpp"
#include "link.hpp"

#define MIKE_QUEUE_LEN 4
#define MIKE_PIN_CLK  0
#define MIKE_PIN_DATA 34
#define MIKE_READ_LEN (MAX_MSG_LEN - 16)
#define MIKE_GAIN_FACTOR 3


struct SampleFreq {
  const static byte Hz_22050 {0x00};
  const static byte Hz_44100 {0x01};
};


class Mike {
    public:
        Mike();
        void init();
        DataUnion* get_msg();
        SemaphoreHandle_t* get_current_semaphore();
        void get_data(byte* data, byte* length);
        static uint32_t sample_rate;
        static bool set_parameters(byte* params, int* length);
    private:
        SemaphoreHandle_t msg_mutex[MIKE_QUEUE_LEN];
        DataUnion msg[MIKE_QUEUE_LEN];
        int current_msg;
        int8_t buffer[MIKE_READ_LEN];
        int16_t *adcBuffer;
};

static TaskHandle_t taskMike;

bool start_Mike_task(Connection* connection);
bool stop_Mike_task();
void Mike_task(void* connection);

#endif