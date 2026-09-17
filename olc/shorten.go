package olc

import (
	"fmt"
	"math"
)

// MinTrimmableCodeLen is the minimum full-code length that can be shortened.
const MinTrimmableCodeLen = 6

var pairResolutions = [...]float64{20.0, 1.0, 0.05, 0.0025, 0.000125}

// Shorten removes leading characters from a full Open Location Code using a reference location.
func Shorten(code string, lat, lng float64) (string, error) {
	dst, err := AppendShorten(nil, code, lat, lng)
	if err != nil {
		return code, err
	}
	return string(dst), nil
}

// AppendShorten appends the shortened form of a full Open Location Code to dst
// using a reference location. It is the zero-allocation variant of Shorten
// when the caller reuses a buffer.
func AppendShorten(dst []byte, code string, lat, lng float64) ([]byte, error) {
	if err := CheckFull(code); err != nil {
		return dst, err
	}
	for i := 0; i < len(code); i++ {
		if code[i] == Padding {
			return dst, ErrPadded
		}
	}

	var upperBuf [MaxEncodedLen]byte
	n := copyUpper(upperBuf[:], code)

	var digits [maxCodeLen]byte
	area := decodeDigits(digits[:stripInto(digits[:], code)])
	if area.Len < MinTrimmableCodeLen {
		return dst, fmt.Errorf("%w: need at least %d", ErrTooShort, MinTrimmableCodeLen)
	}

	lat, lng = clipLatitude(lat), normalizeLng(lng)
	centerLat, centerLng := area.Center()
	distance := math.Max(math.Abs(centerLat-lat), math.Abs(centerLng-lng))

	for i := len(pairResolutions) - 2; i >= 1; i-- {
		if distance < pairResolutions[i]*0.3 {
			return append(dst, upperBuf[(i+1)*2:n]...), nil
		}
	}
	return append(dst, upperBuf[:n]...), nil
}

// RecoverNearest recovers the nearest full Open Location Code from a short code and reference point.
func RecoverNearest(code string, lat, lng float64) (string, error) {
	dst, err := AppendRecoverNearest(nil, code, lat, lng)
	if err != nil {
		return code, err
	}
	return string(dst), nil
}

// AppendRecoverNearest appends the nearest full Open Location Code for a short
// code and reference point to dst. A valid full code is appended uppercased and
// otherwise unchanged. It is the zero-allocation variant of RecoverNearest
// when the caller reuses a buffer.
func AppendRecoverNearest(dst []byte, code string, lat, lng float64) ([]byte, error) {
	if err := CheckFull(code); err == nil {
		n := min(len(code), MaxEncodedLen)
		for i := range n {
			dst = append(dst, upper(code[i]))
		}
		return dst, nil
	}
	if err := CheckShort(code); err != nil {
		return dst, err
	}

	lat, lng = clipLatitude(lat), normalizeLng(lng)

	sep := -1
	for i := 0; i < len(code); i++ {
		if code[i] == Separator {
			sep = i
			break
		}
	}
	padLen := sepPos - sep
	resolution := math.Pow(20, float64(2-(padLen/2)))
	halfRes := resolution / 2

	var merged [MaxEncodedLen]byte
	if _, err := EncodeTo(merged[:], lat, lng, 0); err != nil {
		return dst, err
	}
	write := padLen
	for i := 0; i < len(code); i++ {
		merged[write] = upper(code[i])
		write++
	}

	var digits [maxCodeLen]byte
	dlen := 0
	for i := 0; i < write; i++ {
		b := merged[i]
		if b == Separator || b == Padding {
			continue
		}
		digits[dlen] = b
		dlen++
	}
	area := decodeDigits(digits[:dlen])

	centerLat, centerLng := area.Center()
	if lat+halfRes < centerLat && centerLat-resolution >= -latMax {
		centerLat -= resolution
	} else if lat-halfRes > centerLat && centerLat+resolution <= latMax {
		centerLat += resolution
	}
	if lng+halfRes < centerLng {
		centerLng -= resolution
	} else if lng-halfRes > centerLng {
		centerLng += resolution
	}

	var buf [MaxEncodedLen]byte
	n := integerEncodeTo(buf[:], latitudeAsInteger(centerLat), longitudeAsInteger(centerLng), area.Len)
	return append(dst, buf[:n]...), nil
}

func copyUpper(dst []byte, code string) int {
	n := min(len(code), len(dst))
	for i := range n {
		dst[i] = upper(code[i])
	}
	return n
}
