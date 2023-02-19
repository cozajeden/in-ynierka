from enum import Enum


class Commands(bytearray, Enum):
    ID_OK = b'\x01\x01'
    NEW_ID = b'\x01\x0b'
    START_STREAM_MPU = b'\x05\x01\x01'
    START_STREAM_MIKE = b'\x05\x02\x01'
    STOP_STREAM_MPU = b'\x05\x01\x02'
    STOP_STREAM_MIKE = b'\x05\x02\x02'
    SET = b'\x04'
    MPU = b'\x01'
    MIKE = b'\x02'
    DPS_250 = b'\x00'
    DPS_500 = b'\x01'
    DPS_1000 = b'\x02'
    DPS_2000 = b'\x03'
    G_2 = b'\x00'
    G_4 = b'\x01'
    G_8 = b'\x02'
    G_16 = b'\x03'
    FREQ_22050 = b'\x00'
    FREQ_44100 = b'\x01'