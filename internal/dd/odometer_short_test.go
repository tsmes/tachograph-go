package dd

import "testing"

func TestUnmarshalOdometerShort(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int32
		wantErr bool
	}{
		{
			name:  "zero",
			input: []byte{0x00, 0x00, 0x00},
			want:  0,
		},
		{
			name:  "fourteen",
			input: []byte{0x00, 0x00, 0x0e},
			want:  14,
		},
		{
			name:  "0x010000 = 65536",
			input: []byte{0x01, 0x00, 0x00},
			want:  65536,
		},
		{
			name:  "max non-sentinel",
			input: []byte{0xff, 0xff, 0xfe},
			want:  0xfffffe,
		},
		{
			name:  "all-ones sentinel returns 0",
			input: []byte{0xff, 0xff, 0xff},
			want:  0,
		},
		{
			name:    "too short",
			input:   []byte{0xff, 0xff},
			wantErr: true,
		},
		{
			name:    "empty",
			input:   []byte{},
			wantErr: true,
		},
		{
			name:    "too long",
			input:   []byte{0x00, 0x00, 0x00, 0x00},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := UnmarshalOptions{}
			got, err := opts.UnmarshalOdometerShort(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalOdometerShort() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalOdometerShort() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("UnmarshalOdometerShort() = %d, want %d", got, tt.want)
			}
		})
	}
}
