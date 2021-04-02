package internal

// Specification: https://github.com/alizain/ulid
//
// With heavy inspiration/stolen code from:
// https://github.com/imdario/go-ulid
// https://github.com/oklog/ulid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	crockfordsBase32     string = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	crockfordsBase32Size uint64 = uint64(len(crockfordsBase32))
	// Source: http://www.crockford.com/wrmg/base32.html

	timeLen int = 10
	randLen int = 16
	totLen  int = timeLen + randLen
)

func genULID(timestamp time.Time) (string, error) {
	buf := make([]byte, totLen)

	// Fill in the timestamp (10 bytes, with millisecond precision)
	t := uint64(timestamp.UnixNano() / int64(time.Millisecond))
	binary.LittleEndian.PutUint64(buf[:timeLen], t)

	// Generate some random data and fill the buffer with it (16 bytes)
	size, err := rand.Read(buf[timeLen:])
	if err != nil {
		return "", err
	} else if size != randLen {
		return "", fmt.Errorf("got %d bytes from crypto/rand.Read(), wanted %d bytes.", size, randLen)
	}

	// Encode the full buffer to a string
	tmp := make([]byte, totLen)
	for i := timeLen - 1; i >= 0; i-- {
		mod := t % crockfordsBase32Size
		tmp[i] = crockfordsBase32[mod]
		t = (t - mod) / crockfordsBase32Size
	}
	for i := timeLen; i < len(buf); i++ {
		tmp[i] = crockfordsBase32[uint64(buf[i])%crockfordsBase32Size]
	}
	return string(tmp), nil
}
