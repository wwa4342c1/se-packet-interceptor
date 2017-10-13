package net

import "log"
import "fmt"

type BitBuffer struct {
    bytes []byte
    cur_bit uint32
}

// ReadUBitLong reads bits from the buffer and places them into a uint64.  Max
// supported bits is 64 for now, since we place the read bits into a uint64
// and can't shove more than 64 bits into that.
// TODO: figure out how to deal with odd b.cur_bit values (aka 3)
func (b *BitBuffer) ReadUBitLong(num_bits uint64) uint64 {
    var value, mask uint64

    /* Can't read > 64 bits */
    if num_bits > 64 {
        log.Fatal("ReadUBitLong given too many bits to read.")
    }

    /* If we're asked for an not-multiple-of-8 bits, we need to mask off the
       excess.  Construct the mask here. */
    mask = ^uint64(0)  // Set mask to 0xfff..f by complementing 0x00..0
    mask = mask >> uint64(64 - num_bits)  // Trim first num_bits from the mask

    /* Access each byte and put it in the proper place */
    if num_bits < 8 { num_bits = 8 }  // Make sure that we read at least 1 byte
    for i := 0; i < int(num_bits) / 8 ; i++ {
        value |= uint64(b.bytes[i]) << uint64(i * 8)  // Shift byte to right spot and append
    }

    return (value & mask)
}

func (b *BitBuffer) ReadLong() uint32 {
    return uint32(b.ReadUBitLong(32))
}

func (b *BitBuffer) ReadShort() uint16 {
    return uint16(b.ReadUBitLong(16))
}

func (b *BitBuffer) ReadByte() byte {
    return byte(b.ReadUBitLong(8))
}

func (b *BitBuffer) ReadBit() byte {
    var value byte = b.bytes[b.cur_bit >> 3] >> (b.cur_bit & 7)
    b.cur_bit += 1
    return value & 1
}
