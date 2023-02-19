source venv/bin/activate
cd ../file_decoder/
./decoder < ../server/data/000000_values.csv > ../data_tests/input/test.csv
rm -f ../server/data/000000_values.csv
cd ../data_tests/
python splitter.py input splitted
python printer.py splitted splitted_img
python show_delta_t.py splitted/Mic_0_test.csv