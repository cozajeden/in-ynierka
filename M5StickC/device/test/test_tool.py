from enum import Enum
from functools import partial
from queue import Empty, Queue
from socket import (AF_INET, SHUT_RDWR, SO_REUSEADDR, SOCK_STREAM, SOL_SOCKET,
                    socket, timeout)
from struct import unpack
from threading import Event, Thread
from tkinter import (BOTH, BOTTOM, LEFT, NSEW, RIGHT, TOP, Button, Canvas,
                     Checkbutton, Entry, Frame, IntVar, Label, Listbox, Menu,
                     OptionMenu, Radiobutton, Scale, Scrollbar, Spinbox,
                     StringVar, Text, Tk, X, messagebox)
from typing import Dict, Tuple

from commands import Commands
from figure import Chart
from window import Window


def server(indicators: Dict[str, StringVar], send_queue: Queue, receive_queue: Queue, stop_event: Event, chart_queue: Queue):
    indicators["server_started"].set("Server Started")
    soc = socket(AF_INET, SOCK_STREAM)
    soc.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
    soc.bind(("0.0.0.0", 7123))
    soc.settimeout(1)
    soc.listen()
    while not stop_event.is_set():
        try:
            conn, addr = soc.accept()
            print(f"Connected from {addr}")
            client_thread = Thread(
                target=receiver, args=(indicators, conn, addr, send_queue, receive_queue, stop_event, chart_queue),
                daemon=True
            )
            client_thread.start()
            client_thread.join()
        except timeout:
            pass
    soc.shutdown(SHUT_RDWR)
    soc.close()
    indicators["server_started"].set("Server Stopped")

def sender(indicators: Dict[str, StringVar], conn: socket, send_queue: Queue, stop_event: Event):
    indicators["client_connected"].set("Client Connected")
    while not stop_event.is_set():
        try:
            msg = send_queue.get(timeout=1)
            print(f"Sending {msg}")
            conn.send(msg)
        except Empty:
            pass
    print("sender Closing socket")
    try:
        conn.shutdown(SHUT_RDWR)
        conn.close()
    except:
        pass
    indicators["client_connected"].set("Client Disconnected")

def receiver(indicators: Dict[str, StringVar], conn: socket, addr: Tuple[str, int], send_queue: Queue, receive_queue: Queue, stop_event: Event, chart_queue: Queue):
    sender_task = Thread(
        target=sender, args=(indicators, conn, send_queue, stop_event),
        daemon=True
    )
    manual_mode_task = Thread(
        target=client_handler, args=(indicators, send_queue, receive_queue, stop_event, chart_queue),
        daemon=True
    )
    sender_task.start()
    manual_mode_task.start()
    while not stop_event.is_set():
        try:
            raw = conn.recv(8192)
            if len(raw) < 10: print(f"Received {raw}")
        except timeout:
            continue
        except:
            break
        receive_queue.put(raw)
        if raw == b"":
            break
    print("receiver Closing socket")
    try:
        conn.shutdown(SHUT_RDWR)
        conn.close()
    except:
        pass
    manual_mode_task.join()
    sender_task.join()
    indicators["client_connected"].set("Client Disconnected")

def client_handler(indicators: Dict[str, StringVar], send_queue: Queue, receive_queue: Queue, stop_event: Event, chart_queue: Queue):
    while not stop_event.is_set():
        try:
            raw = receive_queue.get(timeout=1)
        except Empty:
            continue
        if raw[0] == 1:
            print("Sending NEW_ID")
            send_queue.put(Commands.NEW_ID.value)
        elif raw[0] == 2:
            print(f"Received ID: {raw[1]}")
            send_queue.put(Commands.ID_OK.value)
        elif raw[0] == 3 and raw[1] == 2:
            length = unpack("<H", raw[14:16])[0]
            print(length, raw[14:16])
            mic = unpack('<' + 'h'*length, raw[16:length*2 + 16])
            chart_queue.put([("mic", mic)])
        elif raw[0] == 3 and raw[1] == 1:
            length = unpack("<H", raw[14:16])[0]
            data = unpack('<' + 'f'*length*6, raw[16:length*24 + 16])
            accel_x = []
            accel_y = []
            accel_z = []
            gyro_x = []
            gyro_y = []
            gyro_z = []
            for i, value in enumerate(data):
                if i % 6 == 0:
                    accel_x.append(value)
                elif i % 6 == 1:
                    accel_y.append(value)
                elif i % 6 == 2:
                    accel_z.append(value)
                elif i % 6 == 3:
                    gyro_x.append(value)
                elif i % 6 == 4:
                    gyro_y.append(value)
                elif i % 6 == 5:
                    gyro_z.append(value)
            chart_queue.put([
                ("accel_x", accel_x),
                ("accel_y", accel_y),
                ("accel_z", accel_z),
                ("gyro_x", gyro_x),
                ("gyro_y", gyro_y),
                ("gyro_z", gyro_z)
            ])
            
        elif raw[0] == 4:
            if raw[1] == 1:
                print("Received SET OK")
            elif raw[1] == 11:
                print("Received SET FAIL")

        elif raw[0] == 5:
            if raw[1] == 1:
                print("Received ON/OFF OK")
            elif raw[1] == 9:
                print("Received ON/OFF FAIL - MEASUREMENT ALREADY RUNNUNG")
            elif raw[1] == 10:
                print("Received ON/OFF FAIL - MEASUREMENT NOT RUNNING")
            elif raw[1] == 11:
                print("Received ON/OFF FAIL - WRONG COMMAND")


def command_loop(
    indicators: Dict[str, StringVar],
    command_queue: Queue,
    send_queue: Queue,
    receive_queue: Queue,
    stop_event: Event,
    stop_main_loop_event: Event,
    chart_queue: Queue):
    while not stop_main_loop_event.is_set():
        try:
            command = command_queue.get(timeout=1)
        except Empty:
            continue

        if command == "start":
            if stop_event.is_set():
                stop_event.clear()
                Thread(
                    target=server, args=(indicators, send_queue, receive_queue, stop_event, chart_queue),
                    daemon=True
                ).start()

        elif command == "stop":
            if not stop_event.is_set():
                stop_event.set()

        elif command == "start mpu":
            if stop_event.is_set():
                print("Start server first")
            else:
                send_queue.put(Commands.START_STREAM_MPU.value)

        elif command == "stop mpu":
            if stop_event.is_set():
                print("Start server first")
            else:
                send_queue.put(Commands.STOP_STREAM_MPU.value)

        elif command == "start mike":
            if stop_event.is_set():
                print("Start server first")
            else:
                send_queue.put(Commands.START_STREAM_MIKE.value)

        elif command == "stop mike":
            if stop_event.is_set():
                print("Start server first")
            else:
                send_queue.put(Commands.STOP_STREAM_MIKE.value)

        elif command == "set mpu":
            if stop_event.is_set():
                print("Start server first")
            else:
                command = Commands.SET + Commands.MPU
                if indicators['choose_dps_choice'].get() == "250 DPS":
                    command += Commands.DPS_250
                elif indicators['choose_dps_choice'].get() == "500 DPS":
                    command += Commands.DPS_500
                elif indicators['choose_dps_choice'].get() == "1000 DPS":
                    command += Commands.DPS_1000
                elif indicators['choose_dps_choice'].get() == "2000 DPS":
                    command += Commands.DPS_2000
                if indicators['choose_g_choice'].get() == "2 G":
                    command += Commands.G_2
                elif indicators['choose_g_choice'].get() == "4 G":
                    command += Commands.G_4
                elif indicators['choose_g_choice'].get() == "8 G":
                    command += Commands.G_8
                elif indicators['choose_g_choice'].get() == "16 G":
                    command += Commands.G_16
                send_queue.put(command)

        elif command == "set mike":
            if stop_event.is_set():
                print("Start server first")
            else:
                command = Commands.SET + Commands.MIKE
                if indicators['choose_frequency_choice'].get() == "22050 Hz":
                    command += Commands.FREQ_22050
                elif indicators['choose_frequency_choice'].get() == "44100 Hz":
                    command += Commands.FREQ_44100
                send_queue.put(command)

def main():
    root = Window()
    root.title("Test")
    root.geometry("800x600")

    command_queue = Queue()
    send_queue = Queue()
    receive_queue = Queue()
    stop_event = Event()
    stop_main_loop_event = Event()

    indicators = {
        "server_started": StringVar(root, 'Server Stopped'),
        "client_connected": StringVar(root, 'Client Disconnected'),
        "choose_frequency_choice": StringVar(root, "22050 Hz"),
        "choose_dps_choice": StringVar(root, "250 DPS"),
        "choose_g_choice": StringVar(root, "2 G"),
    }

    control_frame = Frame(root)
    control_frame.pack(side=TOP, fill=X)
    control_frame2 = Frame(root)
    control_frame2.pack(side=TOP, fill=X)
    indicator_frame = Frame(root)
    indicator_frame.pack(side=BOTTOM, fill=X)
    
    start_button = Button(control_frame, text="Start Server", command=partial(command_queue.put, "start"))
    start_button.pack(side=LEFT)
    stop_button = Button(control_frame, text="Stop Server", command=partial(command_queue.put, "stop"))
    stop_button.pack(side=LEFT)
    
    start_mpu_button = Button(control_frame, text="Start MPU", command=partial(command_queue.put, "start mpu"))
    start_mpu_button.pack(side=LEFT)
    stop_mpu_button = Button(control_frame, text="Stop MPU", command=partial(command_queue.put, "stop mpu"))
    stop_mpu_button.pack(side=LEFT)
    
    start_mike_button = Button(control_frame, text="Start Mike", command=partial(command_queue.put, "start mike"))
    start_mike_button.pack(side=LEFT)
    stop_mike_button = Button(control_frame, text="Stop Mike", command=partial(command_queue.put, "stop mike"))
    stop_mike_button.pack(side=LEFT)

    choose_dps_menu = OptionMenu(control_frame2, indicators["choose_dps_choice"], "250 DPS", "500 DPS", "1000 DPS", "2000 DPS")
    choose_dps_menu.pack(side=LEFT)

    choose_g_menu = OptionMenu(control_frame2, indicators["choose_g_choice"], "2 G", "4 G", "8 G", "16 G")
    choose_g_menu.pack(side=LEFT)

    set_mpu_button = Button(control_frame2, text="Set MPU", command=partial(command_queue.put, "set mpu"))
    set_mpu_button.pack(side=LEFT)

    choose_frequency_menu = OptionMenu(control_frame2, indicators["choose_frequency_choice"], "22050 Hz", "44100 Hz")
    choose_frequency_menu.pack(side=RIGHT)

    set_mike_button = Button(control_frame2, text="Set Mike", command=partial(command_queue.put, "set mike"))
    set_mike_button.pack(side=RIGHT)

    server_started_label = Label(indicator_frame, relief='solid', borderwidth=2, textvariable=indicators["server_started"], font=("Helvetica bold", 12))
    server_started_label.pack(side=LEFT)
    client_connected_label = Label(indicator_frame, relief='solid', borderwidth=2, textvariable=indicators["client_connected"], font=("Helvetica bold", 12))
    client_connected_label.pack(side=LEFT)

    chart_queue = Queue()
    chart = Chart(root, chart_queue)
    chart.pack(side=TOP, fill=X)
    chart.start()
    # We need to stop chart animation before root.destroy()
    root.add_stop_function(chart.stop)

    command_thread = Thread(
        target=command_loop, args=(indicators, command_queue, send_queue, receive_queue, stop_event, stop_main_loop_event, chart_queue),
    )
    command_thread.start()
    command_queue.put("stop")

    root.mainloop()

    if not stop_event.is_set():
        stop_event.set()
    stop_main_loop_event.set()
    command_thread.join()


if __name__ == "__main__":
    main()
