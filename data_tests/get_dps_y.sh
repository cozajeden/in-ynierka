python splitter.py input splitted

python union.py splitted/ImuDPSY_0_kula.csv splitted/ImuDPSY_1_kula.csv splitted/ImuDPSY_01_kula.csv
python union.py splitted/ImuDPSY_2_kula.csv splitted/ImuDPSY_3_kula.csv splitted/ImuDPSY_23_kula.csv
python union.py splitted/ImuDPSY_01_kula.csv splitted/ImuDPSY_23_kula.csv joined/ImuDPSY_kula.csv

python union.py splitted/ImuDPSX_0_kula.csv splitted/ImuDPSX_1_kula.csv splitted/ImuDPSX_01_kula.csv
python union.py splitted/ImuDPSX_2_kula.csv splitted/ImuDPSX_3_kula.csv splitted/ImuDPSX_23_kula.csv
python union.py splitted/ImuDPSX_01_kula.csv splitted/ImuDPSX_23_kula.csv joined/ImuDPSX_kula.csv

python union.py splitted/ImuDPSZ_0_kula.csv splitted/ImuDPSZ_1_kula.csv splitted/ImuDPSZ_01_kula.csv
python union.py splitted/ImuDPSZ_2_kula.csv splitted/ImuDPSZ_3_kula.csv splitted/ImuDPSZ_23_kula.csv
python union.py splitted/ImuDPSZ_01_kula.csv splitted/ImuDPSZ_23_kula.csv joined/ImuDPSZ_kula.csv

python union.py splitted/ImuGY_0_kula.csv splitted/ImuGY_1_kula.csv splitted/ImuGY_01_kula.csv
python union.py splitted/ImuGY_2_kula.csv splitted/ImuGY_3_kula.csv splitted/ImuGY_23_kula.csv
python union.py splitted/ImuGY_01_kula.csv splitted/ImuGY_23_kula.csv joined/ImuGY_kula.csv

python union.py splitted/ImuGX_0_kula.csv splitted/ImuGX_1_kula.csv splitted/ImuGX_01_kula.csv
python union.py splitted/ImuGX_2_kula.csv splitted/ImuGX_3_kula.csv splitted/ImuGX_23_kula.csv
python union.py splitted/ImuGX_01_kula.csv splitted/ImuGX_23_kula.csv joined/ImuGX_kula.csv

python union.py splitted/ImuGZ_0_kula.csv splitted/ImuGZ_1_kula.csv splitted/ImuGZ_01_kula.csv
python union.py splitted/ImuGZ_2_kula.csv splitted/ImuGZ_3_kula.csv splitted/ImuGZ_23_kula.csv
python union.py splitted/ImuGZ_01_kula.csv splitted/ImuGZ_23_kula.csv joined/ImuGZ_kula.csv

python union.py joined/ImuDPSY_kula.csv joined/ImuDPSX_kula.csv joined/temp.csv
python union.py joined/temp.csv joined/ImuDPSZ_kula.csv joined/ImuDPS_kula.csv

python union.py joined/ImuGY_kula.csv joined/ImuGX_kula.csv joined/temp.csv
python union.py joined/temp.csv joined/ImuGZ_kula.csv joined/ImuG_kula.csv

rm -f joined/temp.csv