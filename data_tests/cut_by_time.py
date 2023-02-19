from datetime import datetime
from os import listdir, mkdir
from os.path import join
from sys import argv
from typing import List, Tuple

exec(f"from {argv[1].split('.')[0]} import FILES, time_format")

import pandas as pd

def main(files: List[Tuple[str, str, str, str]]):
    for start_time, end_time, input_file_path, output_folder, output_file  in files:
        if output_folder not in listdir(): mkdir(output_folder)
        start_time = datetime.strptime(start_time, time_format)
        end_time = datetime.strptime(end_time, time_format)
        data = pd.read_csv(input_file_path, parse_dates=['time'])
        data = data[(data['time'] >= start_time) & (data['time'] < end_time)]
        data.to_csv(join(output_folder, output_file), index=False)


if __name__ == "__main__":
    main(FILES)
