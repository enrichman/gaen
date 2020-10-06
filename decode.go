package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
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

// DecodeFromFile decodes a TemporaryExposureKeyExport binary file
func DecodeFromFile(filename string) ([]*TemporaryExposureKey, error) {
	export, err := UnmarshalExportFile(filename)
	if err != nil {
		return nil, err
	}
	return DecodeExport(export)
}

// UnmarshalExportFile unmarshal a TemporaryExposureKeyExport binary file
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

// DecodeExport decodes a TemporaryExposureKeyExport to a []*TemporaryExposureKey
func DecodeExport(export *export.TemporaryExposureKeyExport) ([]*TemporaryExposureKey, error) {
	teks := make([]*TemporaryExposureKey, 0)

	for _, tek := range export.Keys {
		rollingStartInterval := int(*tek.RollingStartIntervalNumber)
		rollingPeriod := int(*tek.RollingPeriod)

		tek := NewTemporaryExposureKey(tek.KeyData, rollingStartInterval, rollingPeriod)
		err := DecodeTEK(tek)
		if err != nil {
			return nil, err
		}
		teks = append(teks, tek)
	}

	return teks, nil
}

// DecodeTEK decodes a TemporaryExposureKey calculating its Rolling Proximity Identifiers
func DecodeTEK(tek *TemporaryExposureKey) error {
	hkdfReader := hkdf.New(sha256.New, tek.ID, nil, []byte("EN-RPIK"))

	rpiKey := make([]byte, 16)
	if _, err := io.ReadFull(hkdfReader, rpiKey); err != nil {
		return err
	}

	rpis, err := NewRollingProximityIdentifiers(rpiKey, tek.RollingStartInterval, tek.RollingPeriod)
	if err != nil {
		return err
	}
	tek.RPIs = rpis

	return nil
}

// NewRollingProximityIdentifiers returns the Rolling Proximity Identifiers from a key, from the specified starting interval and interval
func NewRollingProximityIdentifiers(rpiKey []byte, rollingStartInterval, rollingPeriod int) ([]*RollingProximityIdentifier, error) {
	rpis := make([]*RollingProximityIdentifier, 0)

	for rp := 0; rp < rollingPeriod; rp++ {
		interval := rp + rollingStartInterval
		newRpi, err := NewRollingProximityIdentifier(rpiKey, interval)
		if err != nil {
			return nil, err
		}
		rpis = append(rpis, newRpi)
	}

	return rpis, nil
}

// NewRollingProximityIdentifier creates a Rolling Proximity Identifier for the specified interval
func NewRollingProximityIdentifier(rpiKey []byte, interval int) (*RollingProximityIdentifier, error) {
	cipher, err := aes.NewCipher(rpiKey)
	if err != nil {
		return nil, err
	}

	rpi := &RollingProximityIdentifier{
		ID:       make([]byte, 16),
		Interval: time.Unix(int64(interval*600), 0),
	}

	cipher.Encrypt(rpi.ID, padInterval(interval))

	return rpi, nil
}

// ID is an alias for an ID made of []byte
type ID []byte

// MarshalJSON is used to override the default marshalJSON
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, id.ToBase64())), nil
}

// ToBase64 returns the string representation of a []byte
func (id ID) ToBase64() string {
	return base64.StdEncoding.EncodeToString(id)
}

// ToHEX returns the hex array representation of a []byte
func (id ID) ToHEX() []string {
	hexArr := make([]string, 0)
	for _, i := range id {
		hex := strconv.FormatUint(uint64(i), 16)
		hexArr = append(hexArr, strings.ToUpper(hex))
	}
	return hexArr
}

// ToInt returns the int array representation of a []byte
func (id ID) ToInt() []int {
	intArr := make([]int, 0)
	for _, i := range id {
		intArr = append(intArr, int(i))
	}
	return intArr
}

// TemporaryExposureKey is the daily tracing key
type TemporaryExposureKey struct {
	ID                   ID `json:"id"`
	RollingStartInterval int
	RollingPeriod        int
	RPIs                 []*RollingProximityIdentifier `json:"rpis"`
}

// NewTemporaryExposureKey returns a Temporary Exposure Key
func NewTemporaryExposureKey(id []byte, rollingStartInterval, rollingPeriod int) *TemporaryExposureKey {
	return &TemporaryExposureKey{
		ID:                   id,
		RollingStartInterval: rollingStartInterval,
		RollingPeriod:        rollingPeriod,
		RPIs:                 make([]*RollingProximityIdentifier, 0),
	}
}

// RollingProximityIdentifier is the bluetooth pseudorandom identifier
type RollingProximityIdentifier struct {
	ID       ID        `json:"id"`
	Interval time.Time `json:"interval"`
}

// padInterval is used to creates the padding array for the specified interval
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
