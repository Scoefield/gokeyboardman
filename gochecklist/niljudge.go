package main

import (
	"fmt"
	"io"
)

type Packet struct {
	PackeyType uint8
	PackeyVersion uint8
	Data *Data
}
type Data struct {
	Stat uint8
	Len uint8
	Buf [8]byte
}
func (p *Packet) UnmarshalBinary(b []byte) error {
	if len(b) < 2 {
		return io.EOF
	}

	p.PackeyType = b[0]
	p.PackeyVersion = b[1]
	// 若长度等于2，那么不会new Data
	if len(b) > 2 {
		p.Data = new(Data)
		// Unmarshal(b[i:], p.Data)
	}
	return nil
}

// bad: 未判断指针是否为nil
func main() {
	packet := new(Packet)
	data := make([]byte, 2)
	if err := packet.UnmarshalBinary(data); err != nil {
		fmt.Println("Failed to unmarshal packet")
		return
	}
	fmt.Printf("Stat: %v\n", packet.Data.Stat)
}