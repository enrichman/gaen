package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gaen/export"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/hkdf"
	"google.golang.org/protobuf/proto"
)

func DecodeFromFile(filename string) error {
	export, err := UnmarshalExportFile(filename)
	if err != nil {
		return err
	}
	teks := DecodeExport(export)

	// print something
	fmt.Printf("\nTEK: [%s] - [%s]\n", base64.StdEncoding.EncodeToString(teks[0].ID), teks[0].IDBytes)
	b, _ := json.MarshalIndent(teks[0].RPIs[0], "", "\t")
	fmt.Printf("\nRPI:\n%v\n", string(b))

	return nil
}

func UnmarshalExportFile(filename string) (*export.TemporaryExposureKeyExport, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	in = in[16:]

	export := &export.TemporaryExposureKeyExport{}
	if err := proto.Unmarshal(in, export); err != nil {
		return nil, err
	}
	return export, nil
}

func DecodeExport(export *export.TemporaryExposureKeyExport) []TemporaryExposureKey {
	teks := make([]TemporaryExposureKey, 0)

	for _, tek := range export.Keys {
		tek := NewTemporaryExposureKeyFromGaenTEK(tek)
		tek.RPIs = DecodeTEK(tek)
		teks = append(teks, tek)
	}

	return teks
}

// DecodeFromTEK ..
func DecodeTEK(tek TemporaryExposureKey) []RollingProximityIdentifier {
	hkdfReader := hkdf.New(sha256.New, tek.ID, nil, []byte("EN-RPIK"))
	rpiKey := make([]byte, 16)
	io.ReadFull(hkdfReader, rpiKey)

	rpis := make([]RollingProximityIdentifier, 0)
	for rp := 0; rp < tek.RollingPeriod; rp++ {
		interval := rp + tek.RollingStartInterval
		newRpi := NewRollingProximityIdentifier(rpiKey, interval)
		rpis = append(rpis, newRpi)
	}

	return rpis
}

type TemporaryExposureKey struct {
	ID                   []byte `json:"id"`
	IDBytes              string `json:"id_bytes"`
	RollingPeriod        int
	RollingStartInterval int
	RPIs                 []RollingProximityIdentifier `json:"rpis"`
}

func NewTemporaryExposureKeyFromGaenTEK(tek *export.TemporaryExposureKey) TemporaryExposureKey {
	return TemporaryExposureKey{
		ID:                   tek.KeyData,
		IDBytes:              encodeToHexString(tek.KeyData),
		RollingPeriod:        int(*tek.RollingPeriod),
		RollingStartInterval: int(*tek.RollingStartIntervalNumber),
		RPIs:                 make([]RollingProximityIdentifier, 0),
	}
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
