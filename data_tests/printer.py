from datetime import datetime, timedelta
from os import listdir, mkdir
from os.path import join
from sys import argv

import matplotlib.dates as mdates
import pandas as pd
from matplotlib import pyplot as plt

in_path = argv[1]
out_path = argv[2]
dateformat = "%H:%M:%S.%f"

unit_mapping = {
    "x-g": "x [g]",
    "y-g": "y [g]",
    "z-g": "z [g]",
    "x-dps": "x [DPS]",
    "y-dps": "y [DPS]",
    "z-dps": "z [DPS]",
    "volume": "amplitude",
    "amplitude": "amplitude",
}

def main():
    if out_path not in listdir(): mkdir(out_path)
    files = listdir(in_path)
    names = [file.split('.')[0] + '.png' for file in files]
    names = [join(out_path, name) for name in names]
    files = [join(in_path, file) for file in files if file.endswith('.csv')]
    for name, file in zip(names, files):
        data = pd.read_csv(file, parse_dates=['time'])
        data.sort_values(by='time', inplace=True)
        fig, ax = plt.subplots()
        labels = []
        for (id, unit), df in data.groupby(["sensor_id", "name"]):
            labels.append(f"{id} {unit_mapping[unit]}")
            df.plot(ax=ax, x="time", y="value", figsize=(15, 10), fontsize=15, alpha=0.4)
        plt.legend(labels)
        ax.set_xticks([],minor=False)
        ax.xaxis.set_minor_locator(mdates.AutoDateLocator())
        ax.xaxis.set_minor_formatter(mdates.DateFormatter(dateformat))
        plt.xticks(rotation=45, minor=True)
        if len(labels) > 15: ax.get_legend().remove()
        plt.tight_layout()
        plt.show()
        #plt.savefig(name, bbox_inches='tight')
        plt.close()

def hist():
    if out_path not in listdir(): mkdir(out_path)
    files = listdir(in_path)
    names = [file.split('.')[0] + '.png' for file in files]
    names = [join(out_path, name) for name in names]
    files = [join(in_path, file) for file in files if file.endswith('.csv')]
    for name, file in zip(names, files):
        data = pd.read_csv(file, parse_dates=['time'])
        data.sort_values(by='time', inplace=True)
        fig, ax = plt.subplots()
        labels = []
        for (id, unit), df in data.groupby(["sensor_id", "name"]):
            labels.append(f"{id} {unit_mapping[unit]}")
            df.sort_values(by='time', inplace=True)
            diff = (df['time'] - df['time'].shift(1))[1:]
            print(diff)
            (diff[diff < timedelta(milliseconds=10)]).apply(lambda x: x.total_seconds()).plot.hist(bins=12, ax=ax, alpha=0.4)
            print((diff[diff < timedelta(milliseconds=10)]).max(), diff.max())
        plt.legend(labels)
        ax.set_xticks([],minor=False)
        ax.xaxis.set_minor_locator(mdates.AutoDateLocator())
        ax.xaxis.set_minor_formatter(mdates.DateFormatter(dateformat))
        plt.xticks(rotation=45, minor=True)
        if len(labels) > 15: ax.get_legend().remove()
        plt.tight_layout()
        plt.savefig(name, bbox_inches='tight')
        plt.close()
            
if __name__ == "__main__":
    main()
