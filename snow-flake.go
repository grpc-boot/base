package base

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ModeWait    = 1
	ModeMaxTime = 2
	ModeError   = 3
)

const (
	defaultEnvKey = `MNUM`
)

const (
	maxIndex   = 0x3fff // 2^14 - 1
	maxMachine = 0xff   // 2^8 - 1
)

// SnowFlake 雪花算法接口，1位0，41位毫秒时间戳，8位机器码，14位递增值
type SnowFlake interface {
	// Id 生成id
	Id() (int64, error)
	// Info 根据id获取信息
	Info(id int64) (timestamp int64, machineId uint8, index int16)
}

// GetMachineId 获取机器Id
type GetMachineId func() (id int64, err error)

// GetMachineIdByIp 根据Ip获取机器Id
func GetMachineIdByIp() GetMachineId {
	return func() (id int64, err error) {
		ip, err := LocalIp()
		if err != nil {
			return 0, err
		}

		return strconv.ParseInt(ip[strings.LastIndexByte(ip, '.')+1:], 10, 8)
	}
}

// GetMachineIdByEnv 根据环境变量获取机器Id
func GetMachineIdByEnv(key string) GetMachineId {
	return func() (id int64, err error) {
		if key == "" {
			key = defaultEnvKey
		}

		return strconv.ParseInt(os.Getenv(key), 10, 9)
	}
}

type work func(current time.Time) (err error)

type snowFlake struct {
	mode           uint8
	lastTimeStamp  int64
	timeStampBegin int64
	index          int16
	machId         int64
	step           work
	mutex          sync.Mutex
}

// NewSFByIp ip方式实例化雪花算法
func NewSFByIp(mode uint8, beginSeconds int64) (sfl SnowFlake, err error) {
	return NewSFByMachineFunc(mode, GetMachineIdByIp(), beginSeconds)
}

// NewSFByMachineFunc GetMachineId方式实例化雪花算法
func NewSFByMachineFunc(mode uint8, machindFunc GetMachineId, beginSeconds int64) (sfl SnowFlake, err error) {
	id, err := machindFunc()
	if err != nil {
		return nil, err
	}

	return NewSF(mode, uint8(id&maxMachine), beginSeconds)
}

// NewSF 实例化雪花算法
func NewSF(mode uint8, id uint8, beginSeconds int64) (sfl SnowFlake, err error) {
	if id > 0xff {
		return nil, ErrMachineId
	}

	sf := &snowFlake{
		mode:           mode,
		lastTimeStamp:  time.Now().UnixNano() / 1e6,
		timeStampBegin: beginSeconds * 1000,
		machId:         int64(id) << 14,
	}

	switch mode {
	case ModeMaxTime:
		sf.step = sf.max
	case ModeError:
		sf.step = sf.err
	default:
		sf.step = sf.wait
	}
	return sf, nil
}

func (sf *snowFlake) Id() (int64, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	err := sf.step(time.Now())
	if err != nil {
		return 0, err
	}

	return ((sf.lastTimeStamp - sf.timeStampBegin) << 22) + sf.machId + int64(sf.index), nil
}

func (sf *snowFlake) Info(id int64) (milliTimestamp int64, machineId uint8, index int16) {
	if id <= sf.machId {
		return
	}

	milliTimestamp = (id >> 22) + sf.timeStampBegin
	machineId = uint8((id >> 14) & maxMachine)
	index = int16(id & maxIndex)
	return
}

func (sf *snowFlake) wait(current time.Time) (err error) {
	curTimeStamp := current.UnixNano() / 1e6

	//时钟回拨等待处理
	for curTimeStamp < sf.lastTimeStamp {
		time.Sleep(time.Millisecond * 5)
		curTimeStamp = time.Now().UnixNano() / 1e6
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > maxIndex {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}

func (sf *snowFlake) max(current time.Time) (err error) {
	curTimeStamp := current.UnixNano() / 1e6

	//时钟回拨使用最大时间
	if curTimeStamp < sf.lastTimeStamp {
		curTimeStamp = sf.lastTimeStamp
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > maxIndex {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}

func (sf *snowFlake) err(current time.Time) (err error) {
	curTimeStamp := current.UnixNano() / 1e6
	//时钟回拨直接抛出异常
	if curTimeStamp < sf.lastTimeStamp {
		return ErrTimeBack
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > maxIndex {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}
