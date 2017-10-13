package main

import (
    "./net"
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

func main() {
    var err error
    var handle *pcap.Handle

    // Parse flags
    flag.Parse()

    // Make the UDP layer aware of our netchannel layer
    // TODO: do we need to remap this based on server port?
    // TODO: what about 27005
    layers.RegisterUDPPortLayerType(layers.UDPPort(27015), net.NetChannelLayerType)

    // Open device
    log.Printf("Staring capture on interface \"%s\".", *iface);
    if *pcap_file != "" {
	    handle, err = pcap.OpenOffline(*pcap_file)
    } else {
	    handle, err = pcap.OpenLive(*iface, int32(*snapshot_len), true, pcap.BlockForever)
    }
    if err != nil {
        log.Fatal("failed to open live stream, ", err)
    }
    defer handle.Close()

    var layer_eth layers.Ethernet
    var layer_ipv4 layers.IPv4
    var layer_udp layers.UDP
    var layer_netchan net.NetChannel

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
	parser.DecodeLayers(packet.Data(), &found_layer_types)
	for _, layer_type := range found_layer_types {
		if layer_type == layers.LayerTypeIPv4 {
			log.Printf("IPV4: %s -> %s", layer_ipv4.SrcIP, layer_ipv4.DstIP)
		}
		if layer_type == layers.LayerTypeUDP {
			log.Printf("UDP: %s -> %s", layer_udp.SrcPort, layer_udp.DstPort)
		}
		if layer_type == net.NetChannelLayerType {
			fmt.Println(layer_netchan.HexDump())
		}
	}
    }
}
