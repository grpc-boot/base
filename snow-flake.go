package base

import (
	"errors"
	"math"
	"os"
	"strconv"
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

var (
	ErrOutOfRange = errors.New("out of range")
	ErrTimeBack   = errors.New("time go back")
	ErrMachineId  = errors.New("illegal machine id")
	ErrLogicId    = errors.New("illegal logic id")
)

// SnowFlake 雪花算法接口，1位0，41位毫秒时间戳，8位机器码，2位业务码，12位递增值
type SnowFlake interface {
	// Id 生成id
	Id(logicId uint8) (int64, error)
	// ZapInfo 根据id获取信息
	Info(id int64) (timestamp int64, machineId uint8, logicId uint8, index int16)
}

// GetMachineId 获取机器Id
type GetMachineId func() (id int64, err error)

// GetMachineIdByIp 根据Ip获取机器Id
func GetMachineIdByIp() GetMachineId {
	return func() (id int64, err error) {
		var (
			id32 uint32
			ip   string
		)

		ip, err = LocalIp()
		if err != nil {
			return 0, err
		}

		id32, err = Ip2Long(ip)
		if err != nil {
			return 0, err
		}
		return int64(id32), nil
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

type work func() (err error)

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
func NewSFByIp(mode uint8, timeStampBegin int64) (sfl SnowFlake, err error) {
	return NewSFByMachineFunc(mode, GetMachineIdByIp(), timeStampBegin)
}

// NewSFByMachineFunc GetMachineId方式实例化雪花算法
func NewSFByMachineFunc(mode uint8, machindFunc GetMachineId, timeStampBegin int64) (sfl SnowFlake, err error) {
	id, err := machindFunc()
	if err != nil {
		return nil, err
	}

	return NewSF(mode, uint8(id&math.MaxUint8), timeStampBegin)
}

// NewSF 实例化雪花算法
func NewSF(mode uint8, id uint8, timeStampBegin int64) (sfl SnowFlake, err error) {
	if id > 0xff {
		return nil, ErrMachineId
	}

	sf := &snowFlake{
		mode:           mode,
		lastTimeStamp:  time.Now().UnixNano() / 1e6,
		timeStampBegin: timeStampBegin,
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

func (sf *snowFlake) Id(logicId uint8) (int64, error) {
	if logicId > 3 {
		return 0, ErrLogicId
	}

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	err := sf.step()
	if err != nil {
		return 0, err
	}

	return ((sf.lastTimeStamp - sf.timeStampBegin) << 22) + sf.machId + (int64(logicId) << 12) + int64(sf.index), nil
}

func (sf *snowFlake) Info(id int64) (timestamp int64, machineId uint8, logicId uint8, index int16) {
	timestamp = (id >> 22) + sf.timeStampBegin
	machineId = uint8((id >> 14) & 0xff)
	logicId = uint8((id >> 12) & 3)
	index = int16(id & 0xfff)
	return
}

func (sf *snowFlake) wait() (err error) {
	curTimeStamp := time.Now().UnixNano() / 1e6

	//时钟回拨等待处理
	for curTimeStamp < sf.lastTimeStamp {
		time.Sleep(time.Millisecond * 5)
		curTimeStamp = time.Now().UnixNano() / 1e6
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > 0xfff {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}

func (sf *snowFlake) max() (err error) {
	curTimeStamp := time.Now().UnixNano() / 1e6

	//时钟回拨使用最大时间
	if curTimeStamp < sf.lastTimeStamp {
		curTimeStamp = sf.lastTimeStamp
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > 0xfff {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}

func (sf *snowFlake) err() (err error) {
	curTimeStamp := time.Now().UnixNano() / 1e6
	//时钟回拨直接抛出异常
	if curTimeStamp < sf.lastTimeStamp {
		return ErrTimeBack
	}

	if curTimeStamp == sf.lastTimeStamp {
		sf.index++
		if sf.index > 0xfff {
			return ErrOutOfRange
		}
	} else {
		sf.index = 0
		sf.lastTimeStamp = curTimeStamp
	}
	return nil
}
