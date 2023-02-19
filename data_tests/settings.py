from datetime import datetime
from enum import Enum
from typing import List, Tuple

Name = str
Value = float
Id = int
FileName = str
DecodedLine = Tuple[Id, datetime, Name, Value]
OriginalNames = List[FileName]


class T(str, Enum):

    @classmethod
    def as_list(cls) -> List[str]:
        return [t.value for t in cls]

class Imu(T):
    x_dps = "x-dps"
    y_dps = "y-dps"
    z_dps = "z-dps"
    x_g = "x-g"
    y_g = "y-g"
    z_g = "z-g"

class Mic(T):
    volume = "volume"

class ImuDPS(T):
    x_dps = "x-dps"
    y_dps = "y-dps"
    z_dps = "z-dps"

class ImuG(T):
    x_g = "x-g"
    y_g = "y-g"
    z_g = "z-g"

class ImuDPSX(T):
    x_dps = "x-dps"

class ImuDPSY(T):
    y_dps = "y-dps"

class ImuDPSZ(T):
    z_dps = "z-dps"

class ImuGX(T):
    x_g = "x-g"

class ImuGY(T):
    y_g = "y-g"

class ImuGZ(T):
    z_g = "z-g"

SENSORS:List[T] = [ImuGX, ImuGY, ImuGZ, ImuDPSX, ImuDPSY, ImuDPSZ, Mic]

def get_sensor(sensor: str) -> T:
    for s in SENSORS:
        if sensor in s.as_list():
            return s.__name__
    raise ValueError("Unknown sensor")
