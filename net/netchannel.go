package net

import (
    "log"
    "strconv"
    "encoding/hex"
    crc32 "hash/crc32"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
)

// Register the netchannel layer
var NetChannelLayerType = gopacket.RegisterLayerType(
	2001,
	gopacket.LayerTypeMetadata{
		Name: "NetChannelLayerType",
		Decoder: gopacket.DecodeFunc(DecodeNetChannel),
	},
)

type NetChannel struct {
    layers.BaseLayer
    payload []byte
    SeqNum uint32
    SeqAckNum uint32
    Flags byte
    Checksum uint16
    RelState byte
    Type byte
}

// DecodeFromBytes extracts a NetChannel from a byte array
func (nc *NetChannel) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	log.Print("Started processing.")
	nc.payload = data
	bytes := BitBuffer{data, 0}
	nc.SeqNum = bytes.ReadLong()
	nc.SeqAckNum = bytes.ReadLong()
	nc.Flags = bytes.ReadByte()

	log.Print("\tFLAGS:")
	if (nc.Flags & PacketFlagReliable != 0) {
	    log.Print("\t\tRELIABLE")
	    return nil
	}
	if (nc.Flags & PacketFlagCompressed != 0) {
	    log.Print("\t\tCOMPRESSED")
	    return nil
	}
	if (nc.Flags & PacketFlagEncrypted != 0) {
	    log.Print("\t\tENCRYPTED")
	    return nil
	}
	if (nc.Flags & PacketFlagChoked != 0) {
	    log.Print("\t\tCHOKED")
	    return nil
	}

        if true {
	    nc.Checksum = bytes.ReadShort()
	    offset := bytes.cur_bit >> 3
	    checksum := DoChecksum(bytes, offset)
	    if checksum != nc.Checksum {
	        log.Print("\tChecksum mismatch, datagram invalid.")
	        return nil
	    }
	}

	nc.RelState = bytes.ReadByte()
	log.Print(bytes.ReadUBitLong(4))

	nc.Type = byte(bytes.ReadUBitLong(NetMsgTypeBits))

	return nil
}

func (nc *NetChannel) LayerType() gopacket.LayerType {
	return NetChannelLayerType
}

func (nc *NetChannel) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypeZero
}

func (nc *NetChannel) CanDecode() gopacket.LayerClass {
	return NetChannelLayerType
}

func (nc *NetChannel) LayerContents() []byte {
	return nil
}

func (nc *NetChannel) LayerPayload() []byte {
	return nil
}

func (nc *NetChannel) Payload() []byte {
	return nil
}

// Hexdump prints a hexdump to stdout
func (nc *NetChannel) HexDump() string {
	return hex.Dump(nc.payload)
}

func (nc *NetChannel) String() string {
	var out string
	out += "SEQ:      " + strconv.FormatInt(int64(nc.SeqNum), 10) + "\n"
	out += "SEQAKC:   " + strconv.FormatInt(int64(nc.SeqAckNum), 10) + "\n"
	out += "FLAGS:    " + strconv.FormatInt(int64(nc.Flags), 10) + "\n"
	out += "CHECKSUM: " + strconv.FormatInt(int64(nc.Checksum), 10) + "\n"
	out += "TYPE:     " + strconv.FormatInt(int64(nc.Type), 10) + "\n"
	return out
}

// NetChannel decoder function.
func DecodeNetChannel(data []byte, p gopacket.PacketBuilder) error {
	// Add a NetChannelLayer onto the list of layers that make up the packet
	netchannel := &NetChannel{}
	err := netchannel.DecodeFromBytes(data, p)
	if err != nil {
	    return err
	}

	p.AddLayer(netchannel)
	p.SetApplicationLayer(netchannel)

	return nil
}

func DoChecksum(buf BitBuffer, offset uint32) uint16 {
    checksum := crc32.ChecksumIEEE(buf.bytes[offset:])
    lower_word := uint16(checksum & 0xffff)
    upper_word := uint16((checksum >> 16) & 0xffff)

    /* The below clusterfuck is brought to you by the fact that Golang removed
       XOR from the usual list of built-in operations.  I have no idea why.
       They moved logical compliment from '~' to '^', which is the XOR operator
       in literally all other languages.  '~' does nothing in golang except
       display a warning that users should use '^' instead.  There was
       literally no reason to get rid of XOR.
       
       Anyways, I needed XOR, so the below is logically equivalent. */
    return (lower_word | upper_word) & ^(lower_word & upper_word)
}
