package main

import (
	"crypto/aes"
	"crypto/sha256"
	"gaen/export"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/hkdf"
)

// DecodeFromTEK ..
func DecodeFromTEK(tek *export.TemporaryExposureKey) []RollingProximityIdentifier {
	rpis := make([]RollingProximityIdentifier, 0)

	hkdfReader := hkdf.New(sha256.New, tek.KeyData, nil, []byte("EN-RPIK"))
	rpiKey := make([]byte, 16)
	io.ReadFull(hkdfReader, rpiKey)

	for rp := 0; rp < int(*tek.RollingPeriod); rp++ {
		interval := rp + int(*tek.RollingStartIntervalNumber)
		newRpi := NewRollingProximityIdentifier(rpiKey, interval)
		rpis = append(rpis, newRpi)
	}

	return rpis
}

// RollingProximityIdentifier ..
type RollingProximityIdentifier struct {
	ID       []byte    `json:"id"`
	IDBytes  string    `json:"id_bytes"`
	Interval time.Time `json:"interval"`
}

// NewRollingProximityIdentifier ..
func NewRollingProximityIdentifier(rpiKey []byte, interval int) RollingProximityIdentifier {
	rpi := RollingProximityIdentifier{
		ID:       make([]byte, 16),
		Interval: time.Unix(int64(interval*600), 0),
	}

	cipher, _ := aes.NewCipher(rpiKey)
	cipher.Encrypt(rpi.ID, padInterval(interval))

	rpi.IDBytes = encodeToHexString(rpi.ID)

	return rpi
}

func encodeToHexString(id []byte) string {
	idBytes := make([]string, 0)
	for _, i := range id {
		hex := strconv.FormatUint(uint64(i), 16)
		idBytes = append(idBytes, strings.ToUpper(hex))
	}
	return strings.Join(idBytes, " ")
}

func padInterval(interval int) []byte {
	// EN-RPI000000
	pad := []byte("EN-RPI")
	for i := 0; i < 6; i++ {
		pad = append(pad, 0)
	}

	pad = append(pad, byte(interval&0xFF))
	pad = append(pad, byte(interval>>8&0xFF))
	pad = append(pad, byte(interval>>16&0xFF))
	pad = append(pad, byte(interval>>24&0xFF))

	return pad
}
