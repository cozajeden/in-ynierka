from datetime import datetime

import pandas as pd

if __name__ == "__main__":
    time_format = "%Y-%m-%d %H:%M:%S.%f"
    impulseTimes = [
        ("2022-10-12 18:29:10.000000", "2022-10-12 18:29:10.200000"),
        ("2022-10-12 18:29:10.800000", "2022-10-12 18:29:11.000000"),
        ("2022-10-12 18:29:11.800000", "2022-10-12 18:29:12.000000"),
        ("2022-10-12 18:29:12.800000", "2022-10-12 18:29:13.000000"),
        ("2022-10-12 18:29:13.800000", "2022-10-12 18:29:14.000000"),
        ("2022-10-12 18:29:14.800000", "2022-10-12 18:29:15.000000"),
        ("2022-10-12 18:29:15.800000", "2022-10-12 18:29:16.000000"),
        ("2022-10-12 18:29:16.800000", "2022-10-12 18:29:17.000000"),
        ("2022-10-12 18:29:17.800000", "2022-10-12 18:29:18.000000"),
        ("2022-10-12 18:29:18.800000", "2022-10-12 18:29:19.000000"),
        ("2022-10-12 18:29:19.800000", "2022-10-12 18:29:20.000000"),
    ]

    inputFilePath = "input/mic.csv"
    data = pd.read_csv(inputFilePath, parse_dates=["time"])
    data.sort_values(["time"], inplace=True)

    impulsesData = []
    for impulseTime in impulseTimes:
        impulseData = data[(data['time'] >= impulseTime[0]) & (data['time'] <= impulseTime[1])]
        impulsesData.append(impulseData)

    sensorIds = range(0, 4, 1)
    for impulseIndex, impulseData in enumerate(impulsesData):
        times = []
        for sensorId in sensorIds:
            time = impulseData[(impulseData['sensor_id'] == sensorId) & (impulseData['value'].abs() >= 10000)].iloc[0].time
            times.append(time)
            print(f'First read over 10000 volume '
                  f'for sensor {sensorId} in '
                  f'impulse nr {impulseIndex}: {time}')
        timesPd = pd.DataFrame(times).sort_values(0)
        greatestDifference = timesPd.iloc[3] - timesPd.iloc[0]
        smallestDifference = timesPd.iloc[2] - timesPd.iloc[1]
        print(f'Greatest difference '
              f'for this impulse '
              f'is {greatestDifference.dt.microseconds}')
        print(f'Smallest difference '
              f'for this impulse '
              f'is {smallestDifference.dt.microseconds}')