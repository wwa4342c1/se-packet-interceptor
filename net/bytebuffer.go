package net

type ByteBuffer struct {
    bytes []byte
    cur_bit uint32
}

func (b *ByteBuffer) ReadUBitLong(num_bits uint32) uint32 {
    var bitmask, start_bit, last_bit, w1_offset, w2_offset, dw1, dw2 uint32

    start_bit = b.cur_bit & uint32(31)
    last_bit = b.cur_bit + num_bits - 1
    w1_offset = b.cur_bit >> 5
    w2_offset = last_bit >> 5
    bitmask = (2 << (num_bits - 1)) - 1

    b.cur_bit += num_bits

    dw1 = uint32(b.bytes[w1_offset]) >> start_bit
    dw2 = uint32(b.bytes[w2_offset]) << (32 - start_bit)

    return (dw1 | dw2) & bitmask
}

func (b *ByteBuffer) ReadLong() uint32 {
    return uint32(b.ReadUBitLong(32))
}

func (b *ByteBuffer) ReadShort() uint16 {
    return uint16(b.ReadUBitLong(16))
}

func (b *ByteBuffer) ReadByte() byte {
    return byte(b.ReadUBitLong(8))
}

func (b *ByteBuffer) ReadBit() byte {
    var value byte = b.bytes[b.cur_bit >> 3] >> (b.cur_bit & 7)
    b.cur_bit += 1
    return value & 1
}
