const { createApp } = Vue
const onCommand = "turn-on"
const offCommand = "turn-off"
const imuSettingsCommand = "imu-settings"
const micSettingsCommand = "mic-settings"
const commandEndpoint = "/command"
const SUCCESS = "success"
const FAILURE = "failure"

createApp({
    data() {
        return {
            sensors: sensors,
            imuG: imuG,
            imuDPS: imuDPS,
            mic: mic,
            messages: [],
            SUCCESS,
            FAILURE,
            chosenId: 0,
            chosenImuG: 0,
            chosenImuDPS: 0,
            chosenMic: 0,
            statusSensor: 1
        }
    },
    methods: {
        sendCommand(payload) {
            that = this
            payload.id = this.chosenId
            axios.post(commandEndpoint, payload)
                .then(function (response) {
                    that.dispatchMessage({
                        kind: SUCCESS,
                        content: response.data.status,
                    })
                })
                .catch(function (error) {
                    that.dispatchMessage({
                        kind: FAILURE,
                        content: error.response.data.error
                    })
                });
        },
        dispatchMessage(message) {
            let today = new Date();
            message.time = String(today.getHours()).padStart(2,"0")
                + ":" + String(today.getMinutes()).padStart(2,"0")
                + ":" + String(today.getSeconds()).padStart(2,"0")
            this.messages.push(message)
        },
        turnOn() {
            this.sendCommand({
                command: onCommand,
                settings: {
                    "sensor_type": parseInt(this.statusSensor)
                }
            })
        },
        turnOff() {
            this.sendCommand({
                command: offCommand,
                settings: {
                    "sensor_type": parseInt(this.statusSensor)
                }
            })
        },
        sendImuConfig() {
            this.sendCommand({
                command: imuSettingsCommand,
                settings: {
                    "sensor_type": 1,
                    "dps": parseInt(this.chosenImuDPS),
                    "g": parseInt(this.chosenImuG)
                }
            })
        },
        sendMicConfig() {
            this.sendCommand({
                command: micSettingsCommand,
                settings: {
                    "sensor_type": 2,
                    "hz": parseInt(this.chosenMic)
                }
            })
        }
    }
}).mount('#app')
