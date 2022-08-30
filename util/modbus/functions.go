package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ReadingName formats MBMD reading names
func ReadingName(val string) string {
	//lint:ignore SA1019 as Title is safe on ascii
	return strings.Title(val)
}

func RTUFloat64ToFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)
	return math.Float64frombits(bits)
}

// decodeMask converts a bit mask in decimal or hex format to uint64
func decodeMask(mask string) (uint64, error) {
	mask = strings.ToLower(mask)

	if strings.HasPrefix(mask, "0x") {
		if len(mask) < 3 {
			return 0, fmt.Errorf("invalid mask: %s", mask)
		}

		b, err := hex.DecodeString(mask[2:])
		if err != nil {
			return 0, fmt.Errorf("invalid mask: %w", err)
		}

		var u uint64
		for _, v := range b {
			u = u<<8 | uint64(v)
		}

		return u, nil
	}

	return strconv.ParseUint(mask, 10, 64)
}

// decodeBool16 converts a masked uint16 to a bool
func decodeBool16(mask uint64) func(b []byte) float64 {
	return func(b []byte) float64 {
		u := binary.BigEndian.Uint16(b)
		if mask > 0 {
			u = u & uint16(mask)
		}
		if u > 0 {
			return 1
		}
		return 0
	}
}

func decodeNaN16(nan uint16, f func(b []byte) float64) func(b []byte) float64 {
	return func(b []byte) float64 {
		if binary.BigEndian.Uint16(b) == nan {
			return 0
		}
		return f(b)
	}
}

func decodeNaN32(nan uint32, f func(b []byte) float64) func(b []byte) float64 {
	return func(b []byte) float64 {
		if binary.BigEndian.Uint32(b) == nan {
			return 0
		}
		return f(b)
	}
}

func decodeNaN64(nan uint64, f func(b []byte) float64) func(b []byte) float64 {
	return func(b []byte) float64 {
		if binary.BigEndian.Uint64(b) == nan {
			return 0
		}
		return f(b)
	}
}
