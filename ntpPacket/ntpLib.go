package ntpPacket

import (
	"encoding/binary"
	"net"
	"time"
)

type ntpPacket []byte

func log2ToDuration(b byte) time.Duration {
	n := int8(b)

	switch {
	case n > 0:
		return time.Duration(uint64(time.Second) << uint(n))
	case n < 0:
		return time.Duration(uint64(time.Second) >> uint(-n))
	default:
		return time.Second
	}
}

func binaryToTime(sec uint64, frac uint64) time.Time {

	// nsec = seconds*10^9  +  (fraction*10^9) / 2^32
	nsec := (sec * 1e9) + ((frac * 1e9) >> 32)

	//loc, _ := time.LoadLocation("Local")
	//now := time.Date(1900, 1, 1, 0, 0, 0, 0, loc).Add(time.Duration(nsec))

	now := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(nsec))

	return now
}

func NewNtpPacket() ntpPacket {
	ret := make(ntpPacket, 48)
	ret.SetLeap(0)
	ret.SetVersion(3)
	ret.SetMode(3)
	return ret
}

func (p ntpPacket) SetLeap(leap uint8) {
	p[0] |= (leap << 6)
}

func (p ntpPacket) SetVersion(version uint8) {
	p[0] |= ((version & 0x07) << 3)
}

func (p ntpPacket) SetMode(mode uint8) {
	p[0] |= (mode & 0x07)
}

func (p ntpPacket) GetLeap() uint8 {
	return uint8((p[0] >> 6) & 0x03) //two first bits
}

func (p ntpPacket) GetVersion() uint8 {
	return uint8((p[0] >> 3) & 0x07)
}

func (p ntpPacket) GetMode() uint8 {
	return uint8(p[0] & 0x07)
}

func (p ntpPacket) Getstratum() uint8 {
	return uint8(p[1])
}

//GetPollInterval translates ntp data in log2 (log base 2) seconds format to
//decimal seconds
func (p ntpPacket) GetPollInterval() time.Duration {
	return log2ToDuration(p[2])
}

func (p ntpPacket) Getprecision() time.Duration {
	return log2ToDuration(p[3])
}

func (p ntpPacket) GetRootDelay() time.Duration { //specs from rfc5905

	secRaw := int16((uint16(p[4]) << 8) | uint16(p[5]))
	fractionRaw := uint16((uint16(p[6]) << 8) | uint16(p[7]))

	sec := int64(secRaw)
	fraction := int64(fractionRaw)

	// nsec = seconds*10^9  +  (fraction*10^9) / 2^32
	nsec := (sec * 1e9) + ((fraction * 1e9) >> 32)

	ret := time.Duration(time.Duration(nsec) * time.Second)

	return ret
}

func (p ntpPacket) GetRootDispersion() time.Duration { //specs from rfc5905

	secRaw := int16((uint16(p[8]) << 8) | uint16(p[9]))
	fractionRaw := uint16((uint16(p[10]) << 8) | uint16(p[11]))

	sec := int64(secRaw)
	fraction := int64(fractionRaw)

	// nsec = seconds*10^9  +  (fraction*10^9) / 2^32
	nsec := (sec * 1e9) + ((fraction * 1e9) >> 32)

	ret := time.Duration(time.Duration(nsec) * time.Second)

	return ret
}

func (p ntpPacket) GetRefClokId() string {
	return string(p[12:16])
}

func (p ntpPacket) GetRefTimestamp() time.Time {

	sec := uint64((uint64(p[16]) << 24) | (uint64(p[17]) << 16) |
		(uint64(p[18]) << 8) | uint64(p[19]))

	fraction := uint64((uint64(p[20]) << 24) | (uint64(p[21]) << 16) |
		(uint64(p[22]) << 8) | uint64(p[23]))

	return binaryToTime(sec, fraction)
}

func (p ntpPacket) GetOriginTimestamp() time.Time {

	sec := uint64((uint64(p[24]) << 24) | (uint64(p[25]) << 16) |
		(uint64(p[26]) << 8) | uint64(p[27]))

	fraction := uint64((uint64(p[28]) << 24) | (uint64(p[29]) << 16) |
		(uint64(p[30]) << 8) | uint64(p[31]))

	return binaryToTime(sec, fraction)
}

func (p ntpPacket) GetRxTimestamp() time.Time {

	sec := uint64((uint64(p[32]) << 24) | (uint64(p[33]) << 16) |
		(uint64(p[34]) << 8) | uint64(p[35]))

	fraction := uint64((uint64(p[36]) << 24) | (uint64(p[37]) << 16) |
		(uint64(p[38]) << 8) | uint64(p[39]))

	return binaryToTime(sec, fraction)
}

func (p ntpPacket) GetTxTimestamp() time.Time {

	sec := uint64((uint64(p[40]) << 24) | (uint64(p[41]) << 16) |
		(uint64(p[42]) << 8) | uint64(p[43]))

	fraction := uint64((uint64(p[44]) << 24) | (uint64(p[45]) << 16) |
		(uint64(p[46]) << 8) | uint64(p[47]))

	return binaryToTime(sec, fraction)
}

func (p ntpPacket) SendTo(socket net.Conn) error {

	err := binary.Write(socket, binary.BigEndian, p)

	if err != nil {
		return err
	}

	return nil
}

func (p ntpPacket) ReadFrom(socket net.Conn) error {

	err := binary.Read(socket, binary.BigEndian, p)

	if err != nil {
		return err
	}

	return nil
}

func (p ntpPacket) GetTime() time.Time {
	return p.GetTxTimestamp()
}
