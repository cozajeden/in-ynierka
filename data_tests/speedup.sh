#!/bin/bash
echo SUPER IMPORTANT! RUN THIS COMMAND AFTER SOURCING VENV 

../file_decoder/decoder < ../server/data/000000_values > input/prezentacja.csv
python splitter.py input splitted
python wav_creator.py
python printer.py input printed
xdg-open sound.wav
