package v5

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadWritePubcompPacket(t *testing.T) {
	tt := []struct {
		testname   string
		pid        PacketID
		code       ReasonCode
		properties *Properties
		want       []byte
	}{
		{
			testname:   "omit properties when code = 0",
			pid:        10,
			code:       CodeSuccess,
			properties: nil,
			want:       []byte{112, 2, 0, 10},
		},
		{
			testname: "code = 0 with properties",
			pid:      10,
			code:     CodeSuccess,
			properties: &Properties{
				ReasonString: []byte("a"),
			},
			want: []byte{112, 8, 0, 10, 0, 4, 0x1f, 0, 1, 'a'},
		}, {
			testname:   "code != 0 with properties",
			pid:        10,
			code:       CodeNotAuthorized,
			properties: &Properties{},
			want:       []byte{112, 4, 0, 10, CodeNotAuthorized, 0},
		},
	}

	for _, v := range tt {
		t.Run(v.testname, func(t *testing.T) {
			a := assert.New(t)
			b := make([]byte, 0, 2048)
			buf := bytes.NewBuffer(b)
			pubcomp := &Pubcomp{
				PacketID:   v.pid,
				Properties: v.properties,
				Code:       v.code,
			}
			err := NewWriter(buf).WriteAndFlush(pubcomp)
			a.Nil(err)
			a.Equal(v.want, buf.Bytes())

			bufr := bytes.NewBuffer(buf.Bytes())
			p, err := NewReader(bufr).ReadPacket()
			a.Nil(err)
			rp := p.(*Pubcomp)

			a.Equal(v.code, rp.Code)
			a.Equal(v.properties, rp.Properties)
			a.Equal(v.pid, rp.PacketID)

		})
	}

}