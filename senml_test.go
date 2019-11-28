package senml_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/mainflux/senml"
	"github.com/stretchr/testify/assert"
)

const (
	cborEncoded = "82ac2169626173652d6e616d6522fb40590000000000002369626173652d756e6974200b24fb404100000000000025fb405900000000000000646e616d650164756e697406fb4062c0000000000007fb4072c0000000000002fb404500000000000005fb4024000000000000ab2169626173652d6e616d6522fb40590000000000002369626173652d756e6974200b25fb405900000000000000666e616d652d310164756e697406fb4062c0000000000007fb4072c0000000000004f505fb4024000000000000"
	jsonEncoded = "5b7b22626e223a22626173652d6e616d65222c226274223a3130302c226275223a22626173652d756e6974222c2262766572223a31312c226276223a33342c226273223a3130302c226e223a226e616d65222c2275223a22756e6974222c2274223a3135302c227574223a3330302c2276223a34322c2273223a31307d2c7b22626e223a22626173652d6e616d65222c226274223a3130302c226275223a22626173652d756e6974222c2262766572223a31312c226273223a3130302c226e223a226e616d652d31222c2275223a22756e6974222c2274223a3135302c227574223a3330302c227662223a747275652c2273223a31307d5d"
	xmlEncoded  = "3c73656e736d6c20786d6c6e733d2275726e3a696574663a706172616d733a786d6c3a6e733a73656e6d6c223e3c73656e6d6c20626e3d22626173652d6e616d65222062743d22313030222062753d22626173652d756e69742220627665723d223131222062763d223334222062733d2231303022206e3d226e616d652220753d22756e69742220743d22313530222075743d223330302220763d2234322220733d223130223e3c2f73656e6d6c3e3c73656e6d6c20626e3d22626173652d6e616d65222062743d22313030222062753d22626173652d756e69742220627665723d223131222062733d2231303022206e3d226e616d652d312220753d22756e69742220743d22313530222075743d22333030222076623d22747275652220733d223130223e3c2f73656e6d6c3e3c2f73656e736d6c3e"
)

var (
	value = 42.0
	sum   = 10.0
	boolV = true
)

var p = senml.Pack{
	Records: []senml.Record{
		senml.Record{
			BaseName:    "base-name",
			BaseTime:    100,
			BaseUnit:    "base-unit",
			BaseVersion: 11,
			BaseSum:     100,
			BaseValue:   34,
			Name:        "name",
			Unit:        "unit",
			Time:        150,
			UpdateTime:  300,
			Value:       &value,
			Sum:         &sum,
		},
		senml.Record{
			BaseName:    "base-name",
			BaseTime:    100,
			BaseUnit:    "base-unit",
			BaseVersion: 11,
			BaseSum:     100,
			Name:        "name-1",
			Unit:        "unit",
			Time:        150,
			UpdateTime:  300,
			BoolValue:   &boolV,
			Sum:         &sum,
		},
	},
}

func TestEncode(t *testing.T) {
	jsnVal, err := hex.DecodeString(jsonEncoded)
	assert.Nil(t, err, "Decoding JSON expected to succeed")
	xmlVal, err := hex.DecodeString(xmlEncoded)
	assert.Nil(t, err, "Decoding XML expected to succeed")
	cborVal, err := hex.DecodeString(cborEncoded)
	assert.Nil(t, err, "Decoding CBOR expected to succeed")
	cases := []struct {
		desc string
		enc  []byte
		p    senml.Pack
		kind senml.Format
		err  error
	}{
		{
			desc: "encode JSON successfully",
			enc:  jsnVal,
			p:    p,
			kind: senml.JSON,
			err:  nil,
		},
		{
			desc: "encode XML successfully",
			enc:  xmlVal,
			p:    p,
			kind: senml.XML,
			err:  nil,
		},
		{
			desc: "encode CBOR successfully",
			enc:  cborVal,
			p:    p,
			kind: senml.CBOR,
			err:  nil,
		},
		{
			desc: "encode unsupported format",
			enc:  nil,
			p:    p,
			kind: 44,
			err:  senml.ErrUnsupportedFormat,
		},
	}
	for _, tc := range cases {
		enc, err := senml.Encode(tc.p, tc.kind)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s expected %s, got %s", tc.desc, tc.err, err))
		assert.Equal(t, tc.enc, enc, fmt.Sprintf("%s expected %s, got %s", tc.desc, tc.enc, enc))
	}
}

func TestDecode(t *testing.T) {
	jsnVal, err := hex.DecodeString(jsonEncoded)
	assert.Nil(t, err, "Decoding JSON expected to succeed")
	xmlVal, err := hex.DecodeString(xmlEncoded)
	assert.Nil(t, err, "Decoding XML expected to succeed")
	cborVal, err := hex.DecodeString(cborEncoded)
	assert.Nil(t, err, "Decoding CBOR expected to succeed")
	xmlPack := p
	xmlPack.Xmlns = "urn:ietf:params:xml:ns:senml"
	cases := []struct {
		desc string
		enc  []byte
		p    senml.Pack
		kind senml.Format
		err  error
	}{
		{
			desc: "encode JSON successfully",
			enc:  jsnVal,
			p:    p,
			kind: senml.JSON,
			err:  nil,
		},
		{
			desc: "encode XML successfully",
			enc:  xmlVal,
			p:    xmlPack,
			kind: senml.XML,
			err:  nil,
		},
		{
			desc: "encode CBOR successfully",
			enc:  cborVal,
			p:    p,
			kind: senml.CBOR,
			err:  nil,
		},
		{
			desc: "encode unsupported format",
			enc:  nil,
			p:    senml.Pack{},
			kind: 44,
			err:  senml.ErrUnsupportedFormat,
		},
	}
	for _, tc := range cases {
		d, err := senml.Decode(tc.enc, tc.kind)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s expected %s, got %s", tc.desc, tc.err, err))
		assert.Equal(t, tc.p, d, fmt.Sprintf("%s expected %v, got %v", tc.desc, tc.p, d))
	}
}

func testValidate(t *testing.T) {

}

func testNormalize(t *testing.T) {

}
