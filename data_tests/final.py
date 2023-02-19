from os import listdir, mkdir
from os.path import join
from sys import argv
from typing import List, Tuple

import pandas as pd


def time_stability():
    pd.set_option("display.precision", 8)
    input_folder = "splitted"
    output_folder = "time_stability"
    if output_folder not in listdir(): mkdir(output_folder)
    imu_files = {
        "title": "IMU",
        1:[
            "ImuDPSY_0_kula.csv",
        ],
        2:[
            "ImuDPSY_0_kula.csv",
            "ImuDPSY_1_kula.csv",
        ],
        3:[
            "ImuDPSY_0_kula.csv",
            "ImuDPSY_1_kula.csv",
            "ImuDPSY_2_kula.csv",
        ],
        4:[
            "ImuDPSY_0_kula.csv",
            "ImuDPSY_1_kula.csv",
            "ImuDPSY_2_kula.csv",
            "ImuDPSY_3_kula.csv",
        ]
    }
    mic_files = {
        "title": "MIC",
        1:[
            "Mic_0_mic.csv",
        ],
        2:[
            "Mic_0_mic.csv",
            "Mic_1_mic.csv",
        ],
        3:[
            "Mic_0_mic.csv",
            "Mic_1_mic.csv",
            "Mic_2_mic.csv",
        ],
        4:[
            "Mic_0_mic.csv",
            "Mic_1_mic.csv",
            "Mic_2_mic.csv",
            "Mic_3_mic.csv",
        ]
    }
    results = []
    for i in range(1,5):
        results.append( f"\nIMU {i} nodes statistics:\n" +
            claculate_time_stability(imu_files[i], input_folder)
        )
        results.append( f"\nMicrophone {i} nodes statistics:\n" +
            claculate_time_stability(mic_files[i], input_folder)
        )
    with open(join(output_folder, "time_stability.txt"), "w") as f:
        f.write("\n".join(results))

def claculate_time_stability(files: List[str], input_folder: str):
    files =  [join(input_folder, f) for f in files]
    files = [pd.read_csv(f, parse_dates=["time"]) for f in files]
    delta_files = []
    for file in files:
        file.sort_values(by="time", inplace=True)
        file["delta_time"] = file["time"].diff()
        file.drop(0, inplace=True)
        delta_files.append(file)
    
    joined = pd.concat(delta_files)
    return str(joined["delta_time"].describe())

def time_sync():
    files = [
        "ImuDPSZ_0_kula.csv",
        "ImuDPSZ_1_kula.csv",
        "ImuDPSZ_2_kula.csv",
        "ImuDPSZ_3_kula.csv",
    ]
    files  = [pd.read_csv(join('splitted', f), parse_dates=["time"]) for f in files]
    maxes  = [file["value"].max()        for file in files]
    mins   = [file["value"].min()        for file in files]
    files  = [file[file!=maxes[0]]       for file in files]
    files  = [file[file!=mins[0]]        for file in files]
    files  = [file[file.value.notnull()] for file in files]
    points = [file.iloc[[0, -1]]         for file in files]
    first_point  = pd.concat([point[:1]  for point in points])
    second_point = pd.concat([point[-1:] for point in points])
    first_diff = []
    second_diff = []
    for i, time1 in enumerate(first_point["time"]):
        first_diff.extend([
            (i, j, abs(time2 - time1).total_seconds()*1_000)
            for j, time2 in enumerate(first_point["time"])
        ])
    for i, time1 in enumerate(second_point["time"]):
        second_diff.extend([
            (i, j, abs(time2 - time1).total_seconds()*1_000)
            for j, time2 in enumerate(second_point["time"])
        ])
    fdf = pd.DataFrame(first_diff, columns=["sensor1", "sensor2", "diff"])
    sdf = pd.DataFrame(second_diff, columns=["sensor1", "sensor2", "diff"])
    fdf = fdf[fdf["sensor1"] != fdf["sensor2"]]
    sdf = sdf[sdf["sensor1"] != sdf["sensor2"]]
    print("Points:")
    print('\n'.join([
        f'sensor: {i}, first: {first_point["time"].iloc[i]}, us: {second_point["time"].iloc[i]}'
        for i in range(4)
    ]))
    print("First point:")
    print('\n'.join([
        f'{sensor1} {sensor2} {diff:0.3f} ms'
        for sensor1, sensor2, diff in first_diff
    ]))
    print("Second point:")
    print('\n'.join([
        f'{sensor1} {sensor2} {diff:0.3f} ms'
        for sensor1, sensor2, diff in second_diff
    ]))
    print("First point statistics:")
    print(fdf.describe())
    print("Second point statistics:")
    print(sdf.describe())
    print("Full statistics:")
    print(pd.concat([fdf,sdf]).describe())


def main():
    time_stability()
    # time_sync()

if __name__ == "__main__":
    main()
