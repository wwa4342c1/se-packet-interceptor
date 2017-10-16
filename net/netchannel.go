package net

import (
    "log"
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

    /* Packet header */
    SeqNum uint32
    SeqAckNum uint32
    Flags byte
    Checksum uint16
    RelState byte

    Cmd byte
}

func (nc *NetChannel) DecodeHeader(bytes *BitBuffer) int {
	nc.SeqNum = bytes.ReadLong()
	//log.Print("\tSEQ: ", nc.SeqNum)
	nc.SeqAckNum = bytes.ReadLong()
	//log.Print("\tSEQACK: ", nc.SeqAckNum)
	nc.Flags = bytes.ReadByte()

	/* Perform checksum */
        nc.Checksum = bytes.ReadShort()
        offset := bytes.cur_bit >> 3
        checksum := bytes.DoChecksum(offset)
        if checksum == nc.Checksum {
            log.Print("\tCHECKSUM.", checksum)
            log.Print("\tChecksum mismatch, datagram invalid.")
            //return nil
        }

	nc.RelState = bytes.ReadByte()
	log.Print("\tRELSTATE: ", nc.RelState)

	if nc.Flags & PacketFlagChoked != 0 {
	    num_choked := bytes.ReadByte()
	    log.Print("\tCHOKED ", num_choked)
	}

	return int(nc.Flags)
}

func (nc *NetChannel) ProcessMessage(bytes *BitBuffer) bool {
	nc.Cmd = byte(bytes.ReadUBitLong(NetMsgTypeBits))
	log.Print("\tCMD: ", nc.Cmd)

	return true
}

// DecodeFromBytes extracts a NetChannel from a byte array
func (nc *NetChannel) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	if data[0] == 0xff && data[1] == 0xff && data[2] == 0xff && data[3] == 0xff {
		// TODO: this also catches OOB queries (aka source server listings), make a thing that catches that and sends them to its own processor
		log.Print("Caught OOB Query, skipping.")
		return nil
	}



	// TODO: add a second handler for incoming/sent packets
	//log.Print("Started processing.")
	nc.payload = data
	//log.Print("\tSIZE:", len(data))
	bytes := BitBuffer{data, 0}

	/* Decode packet header */
	if nc.DecodeHeader(&bytes) == -1 {
		log.Fatal("BAD FLAGS IN PACKET.")
		return nil
	}

	/* Handle reliable packets */
	if nc.Flags & PacketFlagReliable != 0 {
	    log.Print("\t\tRELIABLE")
	    return nil
	}

	if nc.Flags & PacketFlagCompressed != 0 {
	    log.Print("\t\tCOMPRESSED")
	    return nil
	}
	if nc.Flags & PacketFlagEncrypted != 0 {
	    log.Print("\t\tENCRYPTED")
	    return nil
	}

	/* Process message */
	if !nc.ProcessMessage(&bytes) {
		log.Fatal("FAILED TO PROCESS MESSAGE.")
	}

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
