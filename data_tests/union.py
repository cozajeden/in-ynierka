from sys import argv

import pandas as pd

input_file_path1 = argv[1]
input_file_path2 = argv[2]
output_file_path = argv[3]

def main():
    df1 = pd.read_csv(input_file_path1)
    df2 = pd.read_csv(input_file_path2)
    df = pd.concat([df1, df2])
    df.sort_values(by="time").to_csv(output_file_path, index=False)

if __name__ == "__main__":
    main()
