#!/bin/python
# skrypt sprawdza, ktore z obrazow (oprocz svg) sa uzywane w .tex i listuje te, ktore mozna usunac.
# wygodne uzycie: ./image_cleaner.py | xargs -I{} rm {}

from pathlib import Path

picDir = Path("./images")
files = [str(f) for f in picDir.iterdir() if not f.is_dir()]

for name in files:
    with open("./praca_inzynierska.tex") as tex:
        if not name.rsplit("/", 1)[-1].rsplit(".", 1)[0] in tex.read() and name.rsplit(".", 1)[-1] != "svg":
            print(name)
