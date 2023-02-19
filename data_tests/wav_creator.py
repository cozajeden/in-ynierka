import pandas as pd
import wavio

if __name__ == '__main__':
    data = pd.read_csv("splitted/Mic_0_prezentacja.csv", parse_dates=["time"])
    data.sort_values(["time"], inplace=True)
    volumes = data.loc[data['sensor_id'] == 0, ['value']].to_numpy()
    samplerate = 44100
    wavio.write("sound.wav", volumes, samplerate, sampwidth=2)
