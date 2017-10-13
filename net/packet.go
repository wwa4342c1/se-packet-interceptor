package net

const MaxSubChannels int = 8

const (
    SubChannelFree    int = 0  // Subchannel is free to use
    SubChannelToSend  int = 1  // Subchannel has data, but not sent yet
    SubChannelWaiting int = 2  // Subchannel sent data, waiting for ACK
    SubChannelDirty   int = 3  // Subchannel is marked as dirty during changelevel
)

const (
    PacketFlagReliable   byte = (1<<0)  // Packet contains subpacket stream data
    PacketFlagCompressed byte = (1<<1)  // Packet is compressed
    PacketFlagEncrypted  byte = (1<<2)  // Packet is encrypted
    PacketFlagSplit      byte = (1<<3)  // Packet is split
    PacketFlagChoked     byte = (1<<4)  // Packet was choked by sender
    PacketFlagChallenge  byte = (1<<5)  // Packet is a challenge
)

const NetMsgTypeBits uint64 = 8
