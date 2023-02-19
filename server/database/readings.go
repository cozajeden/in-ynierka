package database

import (
	"encoding/binary"
	"fmt"
	"os"
	"server/monitor"
	"server/sensor"
	"sync"
	"time"
)

type ReadingsRepository interface {
	SaveReadings(id uint8, readings []*Reading) error
	Close() error
}

type binaryReadingsRepository struct {
	fileNumber int
	file       *os.File
	saveLock   *sync.RWMutex
}

const gb = 1_000_000_000
const fileSizeLimit = 1 * gb
const fileSavingPath = "/data/"
const databaseTimeFormat = "2006-01-02 15:04:05.000000"
const savingThreads = 5
const waitForDataPeriod = time.Millisecond * 100

var processQueue = make(chan struct {
	readings []*Reading
	id       uint8
}, 1000)
var binaryDataQueue = make(chan []byte, 1000)

func (r *binaryReadingsRepository) SaveReadings(id uint8, readings []*Reading) (err error) {
	processQueue <- struct {
		readings []*Reading
		id       uint8
	}{readings: readings, id: id}

	kind := monitor.IMU
	_, isMic := readings[0].Value[sensor.Volume]
	if isMic {
		kind = monitor.Mic
	}
	monitor.SendDataStruct(len(readings), kind)

	return nil
}

func (r *binaryReadingsRepository) byteifyReadings() {
	pack := <-processQueue
	readings := pack.readings
	id := pack.id
	monitor.ProcessQueueLen(len(processQueue))

	kind := monitor.IMU
	_, isMic := readings[0].Value[sensor.Volume]
	if isMic {
		kind = monitor.Mic
	}
	monitor.ReceiveDataStruct(len(readings), kind)

	start := time.Now()
	binaryData := r.obtainBinary(id, readings)
	monitor.RegisterParseTime(time.Now().Sub(start))
	binaryDataQueue <- binaryData
}

func (r *binaryReadingsRepository) saveBinary() {
	for {
		savingBuffer := 1000
		binaryData := make([]byte, 0, savingBuffer)
		for i := 0; i < savingBuffer; i++ {
			if len(binaryDataQueue) == 0 && len(binaryData) != 0 {
				time.Sleep(waitForDataPeriod)
				if len(binaryDataQueue) == 0 {
					break
				}
			}
			fetchedData := <-binaryDataQueue
			binaryData = append(binaryData, fetchedData...)
		}
		monitor.ByteBufferLen(len(binaryData))
		start := time.Now()
		r.saveLock.Lock()
		err := r.saveToFile(binaryData)
		if err != nil {
			fmt.Printf("couldn't save to file: %v\n", err)
		}
		if r.currentFileSize() > fileSizeLimit {
			err = r.changeFile()
			if err != nil {
				fmt.Printf("couldn't change file: %v\n", err)
			}
		}
		r.saveLock.Unlock()
		monitor.RegisterFileSavingTime(time.Now().Sub(start))
	}
}

func (r *binaryReadingsRepository) Close() error {
	err := r.file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (r *binaryReadingsRepository) currentFileSize() int64 {
	stat, _ := r.file.Stat()
	return stat.Size()
}

func (r *binaryReadingsRepository) changeFile() error {
	err := r.Close()
	if err != nil {
		return err
	}
	err = r.openNewFile()
	if err != nil {
		return err
	}
	return nil
}

func (r *binaryReadingsRepository) obtainBinary(id uint8, readings []*Reading) []byte {
	result := make([]byte, 0, 84*len(readings))
	for _, reading := range readings {
		for readType, value := range reading.Value {
			timeBinary := make([]byte, 8)
			binary.LittleEndian.PutUint64(timeBinary, uint64(reading.Timestamp.UnixMicro()))
			subresult := append([]byte{id, readType}, value...)
			subresult = append(subresult, timeBinary...)
			result = append(result, subresult...)
		}
	}
	return result
}

func (r *binaryReadingsRepository) openNewFile() (err error) {
	r.file, err = os.Create(fileSavingPath + fmt.Sprintf("%06d_values", r.fileNumber))
	if err != nil {
		return err
	}
	r.fileNumber += 1
	return nil
}

func (r *binaryReadingsRepository) saveToFile(data []byte) (err error) {
	_, err = r.file.Write(data)
	if err != nil {
		return fmt.Errorf("couldn't write row %v to file: %v", data, err)
	}
	monitor.WriteBytes(len(data))
	return nil
}

type Reading struct {
	Value     map[uint8][]byte `json:"value"`
	Timestamp time.Time        `json:"timestamp"`
}

func (r *Reading) String() string {
	return fmt.Sprintf("%v: %v", r.Timestamp.Format(databaseTimeFormat), r.Value)
}

func NewReadingsRepository() ReadingsRepository {
	result := &binaryReadingsRepository{
		saveLock: new(sync.RWMutex),
	}
	err := result.openNewFile()
	if err != nil {
		panic(fmt.Sprintf("couldn't open file for db: %v", err))
	}
	for i := 0; i < savingThreads; i++ {
		go func() {
			for {
				result.byteifyReadings()
			}
		}()
	}
	go result.saveBinary()
	return result
}
