package win

import (
	"../bcd"
	"bytes"
	"encoding/binary"
	"fmt"
	//"time"
)

type WinDatetime struct {
	year   uint8
	month  uint8
	day    uint8
	hour   uint8
	minute uint8
	second uint8
}

type WinFormat struct {
	Sequence      uint8
	SubSequence   uint8
	A0            byte
	length        uint16
	datetime      WinDatetime
	year          byte
	month         byte
	day           byte
	hour          byte
	minute        byte
	second        byte
	channel       uint16
	ChannelStatus uint16
	FirstSample   uint32
	Sampling      []int32

	/*
		length       int
		Time         time.Time
		Channel      string
		SamplingSize int
		SamplingRate int
		Buffer       []byte
	*/
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
	binary.Read(buf, binary.BigEndian, &winformat.ChannelStatus)
	binary.Read(buf, binary.BigEndian, &winformat.FirstSample)

	winformat.Sampling = make([]int32, 100)

	fmt.Printf("%0x%0x %X %04X %02d%02d%02d%02d%02d%02d %04X %04X %08X\n", winformat.Sequence,
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
		winformat.ChannelStatus,
		winformat.FirstSample,
	)

	/*
		seq, err := binary.ReadVarint(buf)
		winformat.Sequence = int(seq)

		A0 := buffer[2:3]
		length := buffer[3:5]
		bcddate := buffer[5:11]
		year := bcd.BcdToInt(int(bcddate[0]))
		month := bcd.BcdToInt(int(bcddate[1]))
		day := bcd.BcdToInt(int(bcddate[2]))
		hour := bcd.BcdToInt(int(bcddate[3]))
		minute := bcd.BcdToInt(int(bcddate[4]))
		second := bcd.BcdToInt(int(bcddate[5]))

		ch := buffer[11:13]
		size := buffer[13] >> 4
		rate := (buffer[13]&0x0f)<<8 | buffer[14]&0xff

		firstSample := buffer[15:19]
	*/
	return &winformat
}
