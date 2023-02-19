#ifndef _MPU_
#define _MPU_
#include "globals.hpp"
#include "link.hpp"

#define MPU_ADDRESS           0x68 
#define MPU_WHOAMI            0x75
#define MPU_ACCEL_INTEL_CTRL  0x69
#define MPU_SMPLRT_DIV        0x19
#define MPU_INT_PIN_CFG       0x37
#define MPU_INT_ENABLE        0x38
#define MPU_ACCEL_XOUT_H      0x3B
#define MPU_ACCEL_XOUT_L      0x3C
#define MPU_ACCEL_YOUT_H      0x3D
#define MPU_ACCEL_YOUT_L      0x3E
#define MPU_ACCEL_ZOUT_H      0x3F
#define MPU_ACCEL_ZOUT_L      0x40

#define MPU_TEMP_OUT_H        0x41
#define MPU_TEMP_OUT_L        0x42

#define MPU_GYRO_XOUT_H       0x43
#define MPU_GYRO_XOUT_L       0x44
#define MPU_GYRO_YOUT_H       0x45
#define MPU_GYRO_YOUT_L       0x46
#define MPU_GYRO_ZOUT_H       0x47
#define MPU_GYRO_ZOUT_L       0x48

#define MPU_USER_CTRL         0x6A
#define MPU_PWR_MGMT_1        0x6B
#define MPU_PWR_MGMT_2        0x6C
#define MPU_CONFIG            0x1A
#define MPU_GYRO_CONFIG       0x1B
#define MPU_ACCEL_CONFIG      0x1C
#define MPU_ACCEL_CONFIG2     0x1D
#define MPU_FIFO_EN           0x23

#define MPU_SDA_PIN           21
#define MPU_SCL_PIN           22
#define MPU_INTERRUPT_PIN     35

#define MPU_QUEUE_LEN         4

struct MPUAccelRange {
  const static byte G2 {0x00};
  const static byte G4 {0x01};
  const static byte G8 {0x02};
  const static byte G16 {0x03};
};

struct MPUGyroRange{
  const static byte DPS_250 {0x00};
  const static byte DPS_500 {0x01};
  const static byte DPS_1000 {0x02};
  const static byte DPS_2000 {0x03};
};

class MPU {
    public:
        MPU();
        void init();
        void init(byte a_range, byte g_range);
        void get_data(float* data_ptr);
        void set_accel_range(byte range);
        void set_gyro_range(byte range);
        float get_accel_range();
        float get_gyro_range();
        SemaphoreHandle_t* get_current_semaphore();
        void calculate_from_bytes(uint8_t start_index, float* scale, float* range, float* destination);
        DataUnion* get_msg();
        uint8_t data[14];

        static byte settings[2];
        static bool set_parameters(byte* params, int* length);
        static TwoWire i2c;

    private:
        SemaphoreHandle_t msg_mutex[MPU_QUEUE_LEN];
        DataUnion msg[MPU_QUEUE_LEN];
        int current_msg;

        byte gyro_range;
        byte accel_range;
        float gyro_range_f;
        float accel_range_f;
        float accel_multiplier;
        float gyro_multiplier;

        void reset();
        void set_frequency();
        void write_NBytes(uint8_t driver_Address, uint8_t reg_Address, uint8_t length, uint8_t* data);
        void read_NBytes(uint8_t driver_Address, uint8_t reg_Address, uint8_t length, uint8_t* data);
};

static TaskHandle_t task6886;

bool start_MPU6886_task(Connection *connection);
bool stop_MPU6886_task();
void MPU6886_task(void *connection);

#endif