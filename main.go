package main

import (
    "log"
    "fmt"
    "flag"
    "github.com/google/gopacket"
    "github.com/google/gopacket/pcap"
    "github.com/google/gopacket/layers"
)

var iface = flag.String("iface", "lo", "interface that the game is communicating over")
var pcap_file = flag.String("pcap-file", "", "file to load pcap from")
var snapshot_len = flag.Int("snapshot-length", 1024, "snapshot length.")
var timeout = flag.Int("timeout", 30, "amount of seconds before timeout.")

type NetChannel struct {
    data []byte
}

// Register the netchannel layer
var NetChannelLayerType = gopacket.RegisterLayerType(
	2001,
	gopacket.LayerTypeMetadata{
		"CustomLayerType",
		gopacket.DecodeFunc(decodeNetChannel),
	},
)

// LayerType returns the type of our NetChannel
func (l *NetChannel) LayerType() gopacket.LayerType {
	return NetChannelLayerType
}

// LayerContents returns the contents of our layer.
func (l *NetChannel) LayerContents() []byte {
	return l.data
}

// LayerPayload returns the payload of our layer.
func (l *NetChannel) LayerPayload() []byte {
	return l.data
}

func (l *NetChannel) CanDecode() gopacket.LayerClass {
	return NetChannelLayerType
}

func (l *NetChannel) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypeZero
}

func (l *NetChannel) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

// Hexdump prints a hexdump to stdout
func (l NetChannel) Hexdump() {
	var output string
	for i, v := range l.data {
		if i % 16 == 0 { output += "\n" }
		if i % 8 == 0 { output += "  " }
		output += string(v) + " "
	}
	fmt.Print(output)
}

// NetChannel decoder function.
func decodeNetChannel(data []byte, p gopacket.PacketBuilder) error {
	// Add a NetChannelLayer onto the list of layers that make up the packet
	p.AddLayer(&NetChannel{data})

	// There should be no more bytes after our netchannel thing
	return p.NextDecoder(gopacket.LayerTypeZero)
}

func main() {
    var err error
    var handle *pcap.Handle

    flag.Parse()

    // Make the UDP layer aware of our netchannel layer
    // TODO: do we need to remap this based on server port?
    // TODO: what about 27005
    layers.RegisterUDPPortLayerType(27015, NetChannelLayerType)

    // Open device
    log.Printf("Staring capture on interface \"%s\".", *iface);
    if *pcap_file != "" {
	    handle, err = pcap.OpenOffline(*pcap_file)
    } else {
	    handle, err = pcap.OpenLive(*iface, int32(*snapshot_len), true, pcap.BlockForever)
    }
    if err != nil { log.Fatal("failed to open live stream, ", err) }
    defer handle.Close()

    var layer_eth layers.Ethernet
    var layer_ipv4 layers.IPv4
    var layer_udp layers.UDP
    var layer_netchan NetChannel

    // Decode
    packet_source := gopacket.NewPacketSource(handle, handle.LinkType())
    for packet := range packet_source.Packets() {
        parser := gopacket.NewDecodingLayerParser(
            layers.LayerTypeEthernet,
	    &layer_eth,
	    &layer_ipv4,
	    &layer_udp,
	    &layer_netchan,
	)
	found_layer_types := []gopacket.LayerType{}

	err := parser.DecodeLayers(packet.Data(), &found_layer_types)
	if err != nil {
	    log.Printf("Failed to decode layers: ", err);
	    continue
	}

	for _, layer_type := range found_layer_types {
		if layer_type == layers.LayerTypeIPv4 {
			log.Printf("IPV4: %s -> %s", layer_ipv4.SrcIP, layer_ipv4.DstIP)
		}
		if layer_type == layers.LayerTypeUDP {
			log.Printf("UDP: %s -> %s", layer_udp.SrcPort, layer_udp.DstPort)
		}
		if layer_type == NetChannelLayerType {
			layer_netchan.Hexdump()
		}
	}
    }
}
