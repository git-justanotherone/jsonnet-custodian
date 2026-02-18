package transformers

import (
	"bytes"
	"os"
	"testing"
)

const (
	SOPS_AGE_KEY        = "AGE-SECRET-KEY-1XU0GUPF3E8NELFY9J8K2R6XN4RXGXFTMDRKHYA96FF26DWYW5E8SALHVSG"
	SOPS_DECRYPTED_YAML = `name: test
enabled: true
`
	SOPS_ENCRYPTED_YAML = `name: ENC[AES256_GCM,data:1Z2yJQ==,iv:Fl7iFA9/Zp2+0uicZTMqsW0yrMaycMTz8CZI3Bp31Ng=,tag:Qcpp3XVmIO5HZ/3xKQhrXg==,type:str]
enabled: ENC[AES256_GCM,data:vfjBWQ==,iv:tESUXEqP2zYM9yv9pwaa3E8uSk9gvSFTwxU3NxiO3qM=,tag:miI4IPARN2RhIVk3pKtHSA==,type:bool]
sops:
    age:
        - recipient: age14avqdx2ph6szvewwfnd6yeqrwkpacccz5eg47kzvet9d6586n35qw5esgd
          enc: |
            -----BEGIN AGE ENCRYPTED FILE-----
            YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSB0VDVLa0pZQkhScS9RT3ZK
            bDN4Y3pLNkYyZUlZTml1aXNrUGlFZnp1dENnCkRJRENIT2Z2MlNxOUVnWktNU045
            a2F1VTZMNWlFYkc3YkM1bDNKME9WbmMKLS0tIFJOSm5HRGpzZzFhREFNVEo5Y29h
            QjBZWVJLMlpHU29OZDdueWw4cGZRODQKj7vn3+6Guvos0fA4zKOPpO4VI+BBpJ7d
            C6Gvj8J35T+k70zqstGXg9C/NRqK95HF6cwiu2NwnuthCx4M5R7j+w==
            -----END AGE ENCRYPTED FILE-----
    lastmodified: "2026-01-11T00:59:54Z"
    mac: ENC[AES256_GCM,data:4gcggVLHVoNrMLnvNnlKvA0W57fzIbiOnnfaHKfN8q/11x8hTscxLWnW8N7RO/NIyNFIF8Xwf9f007nbmoWCETk/z6OKFhmHrA6H44UsIYo10cvGBjdRfNe9oc6VIlP9vjnrnhVGcuKLp1KFGdgV/hpf4/Rs2JZU1+Sq2Sf0/r0=,iv:hYDMv/voiR20WLz+TNkL47j7XAoPng32xGxGKZY/LWA=,tag:tr2dPXvdsRaIOK/4Gpwq1A==,type:str]
    unencrypted_suffix: _unencrypted
    version: 3.11.0
`
)

func TestSopsDecryptorTransformer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		foundAt string
		data    []byte
		want    []byte
		wantErr bool
	}{
		struct {
			name    string
			foundAt string
			data    []byte
			want    []byte
			wantErr bool
		}{
			name:    "decrypts sops yaml",
			foundAt: "config.sops.yaml",
			data:    []byte(SOPS_ENCRYPTED_YAML),
			want:    []byte(SOPS_DECRYPTED_YAML),
			wantErr: false,
		},
		{
			name:    "not a sops file",
			foundAt: "random.txt",
			data:    []byte("just some random data"),
			want:    []byte("just some random data"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SOPS_AGE_KEY", SOPS_AGE_KEY)
			got, gotErr := SopsDecryptorTransformer(tt.foundAt, tt.data)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SopsDecryptorTransformer() failed: %v", gotErr)
				}
				return
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("SopsDecryptorTransformer() = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}
