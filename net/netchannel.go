package net

import (
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
    contents []byte
    payload []byte
    butt byte
}

// DecodeFromBytes extracts a NetChannel from a byte array
func (nc *NetChannel) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	nc.payload = data
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
