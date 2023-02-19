from queue import Empty, Queue
from threading import Thread
from tkinter import Tk
from typing import Callable, Tuple

import matplotlib.animation as animation
import matplotlib.pyplot as plt
import numpy as np
from matplotlib.backends.backend_tkagg import FigureCanvasTkAgg
from matplotlib.figure import Figure


class Chart:
    def __init__(self, root: Tk, data_queue: Queue, size: Tuple[int, int]=(9, 9)) -> None:
        self.data = {}
        self.lines = {}
        self.counter = {}
        self.size = size
        self.buffer_size = 4096
        self.data_queue = data_queue
        self.canvas = FigureCanvasTkAgg(self.create_figure(), root)

    @property
    def pack(self) -> Callable[[], None]:
        return self.canvas.get_tk_widget().pack

    def start(self) -> None:
        self.animation = animation.FuncAnimation(self.fig, self.update_plot, interval=120)
        self.update_data_thread = Thread(target=self.update_data, daemon=True)
        self.update_data_thread.start()

    def stop(self) -> None:
        self.animation.event_source.stop()
        self.canvas.get_tk_widget().destroy()
        plt.close(self.fig)

    def create_figure(self) -> Figure:
        self.fig, (self.ax1, self.ax2, self.ax3) = plt.subplots(3, 1, figsize=self.size)
        self.ax1.set_ylim(-3200, 3200)
        self.ax2.set_ylim(-16, 16)
        self.ax3.set_ylim(-2000, 2000)
        self.data['mic'] = np.zeros( self.buffer_size)
        self.lines['mic'] = self.ax1.plot(self.data['mic'], label='mic')[0]
        self.data['accel_x'] = np.zeros( self.buffer_size)
        self.lines['accel_x'] = self.ax2.plot(self.data['accel_x'], label='accel_x')[0]
        self.data['accel_y'] = np.zeros( self.buffer_size)
        self.lines['accel_y'] = self.ax2.plot(self.data['accel_y'], label='accel_y')[0]
        self.data['accel_z'] = np.zeros( self.buffer_size)
        self.lines['accel_z'] = self.ax2.plot(self.data['accel_z'], label='accel_z')[0]
        self.data['gyro_x'] = np.zeros( self.buffer_size)
        self.lines['gyro_x'] = self.ax3.plot(self.data['gyro_x'], label='gyro_x')[0]
        self.data['gyro_y'] = np.zeros( self.buffer_size)
        self.lines['gyro_y'] = self.ax3.plot(self.data['gyro_y'], label='gyro_y')[0]
        self.data['gyro_z'] = np.zeros( self.buffer_size)
        self.lines['gyro_z'] = self.ax3.plot(self.data['gyro_z'], label='gyro_z')[0]
        self.ax1.legend()
        self.ax2.legend()
        self.ax3.legend()
        return self.fig

    def update_data(self):
        while True:
            try:
                data = self.data_queue.get(timeout=0.1)
                for key, value in data:
                    self.data[key] = np.append(self.data[key], value)[- self.buffer_size:]
            except Empty:
                pass

    def update_plot(self, frame_number: int) -> None:
        try:
            for key, value in self.data.items():
                self.lines[key].set_ydata(value)
        except Empty:
            pass
        return self.lines.values()

