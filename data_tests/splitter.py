from os import listdir, mkdir
from os.path import join
from sys import argv
from typing import List, Tuple

import pandas as pd

from settings import OriginalNames, get_sensor

input_folder = argv[1]
output_folder = argv[2]
ext = ".csv"


def get_files() -> Tuple[List[str], OriginalNames]:
    return ([join(input_folder, file) for file in listdir(input_folder) if file.endswith(ext)],
        [file for file in listdir(input_folder) if file.endswith(ext)])

def main():
    if not output_folder in listdir(): mkdir(output_folder) 
    for file, name in zip(*get_files()):
        file = pd.read_csv(file)
        file["sensor"] = file["name"].apply(get_sensor)
        for (sensor, id), file in file.groupby(["sensor", "sensor_id"]):
            file.sort_values(by="time").to_csv(join(output_folder, f"{sensor}_{id}_{name}"), index=False)


if __name__ == "__main__":
    main()
