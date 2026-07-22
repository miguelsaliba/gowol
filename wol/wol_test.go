package wol

import (
	"errors"
	"testing"
)

func TestParseMac(t *testing.T) {
	cases := []struct {
		Mac    string
		Parsed MACAddress
		Error  error
	}{
		{Mac: "11:22:33:44:55:66", Parsed: MACAddress{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, Error: nil},
		{Mac: "FF:FF:FF:FF:FF:FF", Parsed: MACAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, Error: nil},
		{Mac: "FF-FF-FF-FF-FF-FF", Parsed: MACAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, Error: nil},
		{Mac: "111111111111", Parsed: MACAddress{0x11, 0x11, 0x11, 0x11, 0x11, 0x11}, Error: nil},
		{Mac: "11:22:33:44:55:66:77:88", Parsed: MACAddress{}, Error: ErrInvalidMAC},
		{Mac: "11:22", Parsed: MACAddress{}, Error: ErrInvalidMAC},
		{Mac: "11:22:33:44", Parsed: MACAddress{}, Error: ErrInvalidMAC},
	}
	for _, test := range cases {
		t.Run(test.Mac, func(t *testing.T) {
			got, err := ParseMac(test.Mac)

			switch {
			case test.Error == nil && err != nil:
				t.Errorf("unexpected error %v", err)

			case test.Error != nil && err == nil:
				t.Errorf("expected error %v but got nil", test.Error)

			case !errors.Is(err, test.Error):
				t.Errorf("expected error %v but got %v", test.Error, err)

			case test.Parsed != got:
				t.Errorf("expected %v but got %v", test.Parsed, got)
			}
		})

	}
}
