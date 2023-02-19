from tkinter import Tk

class Window(Tk):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.stop_functions = []

    def add_stop_function(self, func):
        self.stop_functions.append(func)

    def destroy(self) -> None:
        for fun in self.stop_functions:
            fun()
        return super().destroy()