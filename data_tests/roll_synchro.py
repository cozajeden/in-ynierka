from datetime import datetime

import pandas as pd

if __name__ == "__main__":
    time_format = "%Y-%m-%d %H:%M:%S.%f"
    edgeTimes = [
        ("2022-10-04 20:37:36.400000", "2022-10-04 20:37:36.500000"),
        ("2022-10-04 20:37:38.400000", "2022-10-04 20:37:38.500000"),
        ("2022-10-04 20:37:40.820000", "2022-10-04 20:37:40.900000"),
        ("2022-10-04 20:37:42.000000", "2022-10-04 20:37:42.200000"),
        ("2022-10-04 20:37:43.300000", "2022-10-04 20:37:43.500000"),
        ("2022-10-04 20:37:44.200000", "2022-10-04 20:37:44.400000"),
        ("2022-10-04 20:37:45.600000", "2022-10-04 20:37:45.800000"),
        ("2022-10-04 20:37:46.700000", "2022-10-04 20:37:46.900000"),
        ("2022-10-04 20:37:48.000000", "2022-10-04 20:37:48.100000"),
        ("2022-10-04 20:37:49.100000", "2022-10-04 20:37:49.300000"),
        ("2022-10-04 20:37:50.200000", "2022-10-04 20:37:50.400000"),
    ]
    inputFilePath = "splitted/mergedY.csv"
    data = pd.read_csv(inputFilePath, parse_dates=["time"])
    data.sort_values(["time"], inplace=True)

    edgesData = []
    for edgeTime in edgeTimes:
        edgeData = data[(data['time'] >= edgeTime[0]) & (data['time'] <= edgeTime[1])]
        edgesData.append(edgeData)

    sensorIds = range(0, 4, 1)
    for edgeIndex, edgeData in enumerate(edgesData):
        times = []
        for sensorId in sensorIds:
            time = edgeData[(edgeData['sensor_id'] == sensorId) &
                            (
                                    (edgeData['value'] == 247.90308) |
                                    (edgeData['value'] == -249.99237)
                            )].iloc[0].time
            times.append(time)
            print(f'First stauration '
                  f'for sensor {sensorId} in '
                  f'edge nr {edgeIndex}: {time}')
        timesPd = pd.DataFrame(times).sort_values(0)
        greatestDifference = timesPd.iloc[3] - timesPd.iloc[0]
        smallestDifference = timesPd.iloc[2] - timesPd.iloc[1]
        print(f'Greatest difference '
              f'for this edge '
              f'is {greatestDifference.dt.microseconds}')
        print(f'Smallest difference '
              f'for this edge '
              f'is {smallestDifference.dt.microseconds}')