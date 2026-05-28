package vu

import (
	"testing"

	"github.com/way-platform/tachograph-go/internal/dd"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// TestParseOneCalibrationRecordGen2V2_CalibrationCountry guards against the same
// enum-mapping regression as the vehicle registration nation: the calibration
// country byte (offset 247) must be mapped through the protocol_enum_value
// annotation, not cast directly to the proto enum. With the raw cast, a
// Norwegian calibration (protocol value 37) decoded as MONACO (proto enum value
// 37). See tacho#1090.
//
// A 252-byte all-zero record parses cleanly (lenient sub-parsers); only the
// calibration country byte is varied.
func TestParseOneCalibrationRecordGen2V2_CalibrationCountry(t *testing.T) {
	const lenRecord = 252
	const idxCalCountry = 247

	tests := []struct {
		name        string
		countryByte byte
		want        ddv1.NationNumeric
	}{
		{name: "Norway", countryByte: 0x25, want: ddv1.NationNumeric_NORWAY},
		{name: "Monaco", countryByte: 0x22, want: ddv1.NationNumeric_MONACO},
		{name: "Default (0x00)", countryByte: 0x00, want: ddv1.NationNumeric_NATION_NUMERIC_DEFAULT},
		{name: "Unrecognized (100)", countryByte: 0x64, want: ddv1.NationNumeric_NATION_NUMERIC_UNRECOGNIZED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, lenRecord)
			data[idxCalCountry] = tt.countryByte

			var opts dd.UnmarshalOptions
			rec, err := parseOneCalibrationRecordGen2V2(opts, data)
			if err != nil {
				t.Fatalf("parseOneCalibrationRecordGen2V2() unexpected error: %v", err)
			}
			if rec.GetCalibrationCountry() != tt.want {
				t.Errorf("GetCalibrationCountry() = %v, want %v", rec.GetCalibrationCountry(), tt.want)
			}
		})
	}
}

// TestParseOneCalibrationRecordGen2V2_OdometerSentinel covers the Continental
// Gen2 V2 ACTIVATION case where newOdometer (offset 152) contains the
// OdometerShort "value not available" sentinel (0xFFFFFF). The parser must
// recognise the sentinel and emit 0, not 16777215. See tacho#1128.
func TestParseOneCalibrationRecordGen2V2_OdometerSentinel(t *testing.T) {
	const lenRecord = 252
	const idxOldOdometer = 149
	const idxNewOdometer = 152

	tests := []struct {
		name    string
		old     []byte
		new     []byte
		wantOld int32
		wantNew int32
	}{
		{
			name:    "normal values",
			old:     []byte{0x00, 0x00, 0x0e}, // 14
			new:     []byte{0x00, 0x00, 0x14}, // 20
			wantOld: 14,
			wantNew: 20,
		},
		{
			name:    "newOdometer sentinel (Continental ACTIVATION)",
			old:     []byte{0x00, 0x00, 0x0e}, // 14
			new:     []byte{0xff, 0xff, 0xff}, // sentinel
			wantOld: 14,
			wantNew: 0,
		},
		{
			name:    "both sentinels",
			old:     []byte{0xff, 0xff, 0xff},
			new:     []byte{0xff, 0xff, 0xff},
			wantOld: 0,
			wantNew: 0,
		},
		{
			name:    "zero is not the sentinel",
			old:     []byte{0x00, 0x00, 0x00},
			new:     []byte{0x00, 0x00, 0x00},
			wantOld: 0,
			wantNew: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, lenRecord)
			copy(data[idxOldOdometer:idxOldOdometer+3], tt.old)
			copy(data[idxNewOdometer:idxNewOdometer+3], tt.new)

			var opts dd.UnmarshalOptions
			rec, err := parseOneCalibrationRecordGen2V2(opts, data)
			if err != nil {
				t.Fatalf("parseOneCalibrationRecordGen2V2() unexpected error: %v", err)
			}
			if got := rec.GetOldOdometerValueKm(); got != tt.wantOld {
				t.Errorf("GetOldOdometerValueKm() = %d, want %d", got, tt.wantOld)
			}
			if got := rec.GetNewOdometerValueKm(); got != tt.wantNew {
				t.Errorf("GetNewOdometerValueKm() = %d, want %d", got, tt.wantNew)
			}
		})
	}
}
