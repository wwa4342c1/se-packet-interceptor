package net

import (
    "log"
    "strconv"
    "encoding/hex"
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
	nc.Checksum = bytes.ReadShort()
	nc.RelState = bytes.ReadByte()
	log.Print(bytes.ReadUBitLong(4))

	if nc.Flags & PacketFlagChoked != 0 {
            //TODO: handle this (in a different layer?)
	    log.Print("Choked.")
	    return nil
	} else if nc.Flags & PacketFlagChallenge != 0 {
            //TODO: handle this (in a different layer?)
	    log.Print("Challenge.")
	    return nil
	} else if nc.Flags & PacketFlagReliable != 0 {
            //TODO: handle this (in a different layer?)
	    log.Print("Reliable.")
	    return nil
	}

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
