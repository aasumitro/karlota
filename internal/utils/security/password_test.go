package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"karlota.aasumitro.id/internal/utils/security"
)

// Mock rand.Reader that always returns an error
type mockRandReader struct {
	err error
}

func (r *mockRandReader) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

//func TestHash_Must_Success(t *testing.T) {
//	type args struct {
//		s string
//	}
//
//	tests := []struct {
//		name    string
//		args    args
//		want    bool
//		wantErr assert.ErrorAssertionFunc
//	}{
//		{
//			name:    "test hash make and verify",
//			args:    args{s: "secret"},
//			want:    true,
//			wantErr: assert.NoError,
//		},
//		{
//			name:    "test hash make and verify",
//			args:    args{s: "12345"},
//			want:    true,
//			wantErr: assert.NoError,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := security.PasswordHash{Stored: "", Supplied: tt.args.s}
//			pwd, err := h.MakePassword(security.Parallelization)
//			h.Stored = pwd
//			assert.Nil(t, err)
//			valid, err := h.ComparePassword(security.Parallelization)
//			assert.Nil(t, err)
//			assert.Equalf(t, tt.want, valid, "Hash(%v)", tt.args.s)
//		})
//	}
//}
//
//func TestHash_Make_Error(t *testing.T) {
//	t.Run("ERROR FROM READER", func(t *testing.T) {
//		expectedErr := errors.New("random read error")
//		mockRand := &mockRandReader{err: expectedErr}
//
//		origRandReader := rand.Reader
//		defer func() { rand.Reader = origRandReader }()
//		rand.Reader = mockRand
//
//		h := security.PasswordHash{Supplied: "password123"}
//		secret, err := h.MakePassword(security.Parallelization)
//		assert.NotNil(t, err)
//		assert.Empty(t, secret)
//		assert.Equal(t, err.Error(), expectedErr.Error())
//	})
//	t.Run("ERROR FROM HASH", func(t *testing.T) {
//		h := security.PasswordHash{Supplied: "password123"}
//		secret, err := h.MakePassword(-1)
//		assert.NotNil(t, err)
//		assert.Empty(t, secret)
//	})
//}
//
//func TestHash_Compare_Error(t *testing.T) {
//	t.Run("ERROR WHEN VALIDATE LEN", func(t *testing.T) {
//		h := security.PasswordHash{Stored: "", Supplied: ""}
//		valid, err := h.ComparePassword(security.Parallelization)
//		assert.NotNil(t, err)
//		assert.False(t, valid)
//	})
//
//	t.Run("ERROR WHEN HEX DECODE", func(t *testing.T) {
//		h := security.PasswordHash{Stored: "1234.gibberish", Supplied: ""}
//		valid, err := h.ComparePassword(security.Parallelization)
//		assert.NotNil(t, err)
//		assert.False(t, valid)
//		assert.Equal(t, err.Error(), security.ErrorPasswordUnableToVerify.Error())
//	})
//
//	t.Run("ERROR WHEN HEX DECODE", func(t *testing.T) {
//		h := security.PasswordHash{Stored: "1234.1234", Supplied: ""}
//		valid, err := h.ComparePassword(-1)
//		assert.NotNil(t, err)
//		assert.False(t, valid)
//		assert.Equal(t, err.Error(), security.ErrorPasswordUnableToVerify.Error())
//	})
//}

func TestHash_PasswordHashArgon2_Success(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "test hash make and verify",
			args:    args{s: "secret"},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name:    "test hash make and verify",
			args:    args{s: "12345"},
			want:    true,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := security.PasswordHashArgon2{Stored: "", Supplied: tt.args.s}
			pwd, err := h.MakePassword()
			h.Stored = pwd
			assert.Nil(t, err)
			valid, err := h.ComparePassword()
			assert.Nil(t, err)
			assert.Equalf(t, tt.want, valid, "Hash(%v)", tt.args.s)
		})
	}

	// Test@1234
	// argon2 algo 57.340366ms to compare - ka4e15tCd6iiS3iJZPH/XwDWy8JqMvTj+AXxRafwOBU=.KE8Z75HYe3wgquy82NxdglWRhjriMZTIB2LjAqRTU/Y=
	// scrypt algo 661.841014ms to compare- bee27f92033b49dbb49cffd7012a0e0129724e9c03f2317a2271455a7c19a1a4.be2b932968c024749e62605c376310eab86fa40f54b9dadead463ff06152dad0
}
