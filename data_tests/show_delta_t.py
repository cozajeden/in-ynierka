from sys import argv

import pandas as pd
from matplotlib import pyplot as plt

in_file = argv[1]

df = pd.read_csv(in_file, parse_dates=['time']).sort_values(by='time')
df['time'].diff()[1:].plot()
print(df["time"].diff()[1:].describe())
plt.show()
