package win

import (
	"../bcd"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type WinFormat struct {
	Sequence    uint8
	SubSequence uint8
	A0          byte
	length      uint16
	year        byte
	month       byte
	day         byte
	hour        byte
	minute      byte
	second      byte
	channel     uint16
	rate        byte
	size        byte
	FirstSample int32
	Sampling    []int32
}

func Parse(buffer []byte) *WinFormat {
	winformat := WinFormat{}
	buf := bytes.NewBuffer(buffer)
	binary.Read(buf, binary.BigEndian, &winformat.Sequence)
	binary.Read(buf, binary.BigEndian, &winformat.SubSequence)
	binary.Read(buf, binary.BigEndian, &winformat.A0)
	binary.Read(buf, binary.BigEndian, &winformat.length)
	binary.Read(buf, binary.BigEndian, &winformat.year)
	binary.Read(buf, binary.BigEndian, &winformat.month)
	binary.Read(buf, binary.BigEndian, &winformat.day)
	binary.Read(buf, binary.BigEndian, &winformat.hour)
	binary.Read(buf, binary.BigEndian, &winformat.minute)
	binary.Read(buf, binary.BigEndian, &winformat.second)
	binary.Read(buf, binary.BigEndian, &winformat.channel)
	binary.Read(buf, binary.BigEndian, &winformat.size)
	binary.Read(buf, binary.BigEndian, &winformat.rate)

	binary.Read(buf, binary.BigEndian, &winformat.FirstSample)

	size := winformat.size >> 4
	rate := (winformat.size&0x0f)<<8 | winformat.rate&0xff

	winformat.Sampling = make([]int32, rate-1)

	fmt.Printf("%0X%0X %X %04X %02d%02d%02d%02d%02d%02d %04X %02X%02X(%d %d) %08X\n", winformat.Sequence,
		winformat.SubSequence,
		int(winformat.A0),
		winformat.length,
		bcd.BcdToInt(int(winformat.year)),
		bcd.BcdToInt(int(winformat.month)),
		bcd.BcdToInt(int(winformat.day)),
		bcd.BcdToInt(int(winformat.hour)),
		bcd.BcdToInt(int(winformat.minute)),
		bcd.BcdToInt(int(winformat.second)),
		winformat.channel,
		winformat.size,
		winformat.rate,
		size,
		rate,
		winformat.FirstSample,
	)

	if size == 0 {
		rate = rate / 2

	}

	for i := 0; i < int(rate)-1; i++ {
		if size == 4 {
			var value int32
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = value
		}
		if size == 3 {
			value := make([]byte, 3)
			binary.Read(buf, binary.BigEndian, &value)
			num, _ := strconv.ParseInt(fmt.Sprintf("%x%x%x", value[0], value[1], value[2]), 16, 32)
			winformat.Sampling[i] = int32(num)
		}
		if size == 2 {
			var value int16
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)

		}
		if size == 1 {
			var value int8
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)

		}
		if size == 0 {
			var value int8
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)
			//fmt.Printf("%d %x %x%x \n", i, value, winformat.Sampling[i]&0xf0, (winformat.Sampling[i])&0x0f)
		}
	}

	for i := 0; i < int(rate)-1; i++ {
		fmt.Printf("%d %x\n", i, winformat.Sampling[i])
	}
	return &winformat
}
