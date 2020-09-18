package j8a

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/hako/durafmt"
	"testing"
	"time"
)

func TestPDurationAsString(t *testing.T) {
	pd399 := PDuration(time.Hour * 24 * 399)
	got := pd399.AsString()
	want := "1 year 4 weeks"
	if got != want {
		t.Errorf("durafmt parse error, want %s, got %s", want, got)
	}
}

func TestPDuration_AsDays(t *testing.T) {
	pd399 := PDuration(time.Hour * 24 * 399)
	got := pd399.AsDays()
	want := 399
	if got != want {
		t.Errorf("durafmt parse error, want %d, got %d", want, got)
	}
}

func TestBrowserExpiry_AsDays(t *testing.T) {
	got := new(TlsLink).browserExpiry().AsDays()
	want := 398
	if got != want {
		t.Errorf("wrong browser expiry days, want %d, got %d", want, got)
	}
}

func TestParseTlsLinks(t *testing.T) {
	tlsConfig := mockTlsConfig()
	c, _ := x509.ParseCertificate(tlsConfig.Certificates[0].Certificate[0])
	tlsLinks := parseTlsLinks([]*x509.Certificate{c})

	logCertStats(tlsLinks)

	if len(tlsLinks) != 1 {
		t.Errorf("tls links parsed incorrectly")
	} else {
		if tlsLinks[0].isCA != false {
			t.Errorf("cert should not be a CA")
		}
		if tlsLinks[0].totalValidity.AsDuration().Seconds() != time.Duration(time.Second*352257299).Seconds() {
			t.Errorf("total validity should be %s", durafmt.Parse(tlsLinks[0].totalValidity.AsDuration()))
		}
	}
}

func TestCheckCertChain(t *testing.T) {
	tlsConfig := mockTlsConfig()
	verified, err := checkCertChain(tlsConfig.Certificates[0])
	if err != nil {
		t.Errorf("certificate chain with 1 TLS cert, 1 root cert not validated, cause: %s", err)
	} else {
		t.Logf("normal. certificate chain with 1 TLS cert, 1 root cert validated, length: %d", len(verified))
	}
}

func TestTlsHealthCheck(t *testing.T) {
	//this only needs to be covered for no runtime exceptions as it logs to console. no assertions.
	tlsHealthCheck(mockTlsConfig(), false)
}

func TestCertChainC_I_R(t *testing.T) {
	certPem := "-----BEGIN CERTIFICATE-----\nMIIFkDCCA3igAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwdjELMAkGA1UEBhMCQVUx\nDDAKBgNVBAgMA05TVzENMAsGA1UECgwEbXljYTEaMBgGA1UECwwRbXljYSBpbnRl\ncm1lZGlhdGUxLjAsBgNVBAMMJW15IGNlcnRpZmljYXRlIGF1dGhvcml0eSBpbnRl\ncm1lZGlhdGUwIBcNMjAwOTE3MjEzMTQ5WhgPMjEwMjExMDcyMTMxNDlaMF8xCzAJ\nBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTENMAsGA1UE\nCgwEY2VydDEQMA4GA1UECwwHY2VydCBvdTEQMA4GA1UEAwwHY2VydCBjbjCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN3FFHDc3fWIyxukMDRriEbYtVA4\n1EeiQiwf7RLdDxh+N2VAazUbbxUJ06nKAslX2+6ZmJrMlS+ionX1BvPhPy3snuZI\n1movXcvH6ZV5yUGZyJDocjOTHHqNwPSDOAQX87tLjQbCa8Rw//B488GoPbaZlWYD\nvZQ0Mw5rasiu0B+OI6PL8+Vnc2jXdPlc3tiNoIVXRZ14TNei7bUDA3O1y593ift2\ntQ/TZxlY7fylZWhTV4sUm/9yk/zob+dyzro795Jy8vThlePAN//tZGLWFzG7a8o9\nMx36BPncSZ0v+EfEvP24ZffIDFRtysBewu2+33IVpISlbaHgj6nsuv8GFM0CAwEA\nAaOCATswggE3MAkGA1UdEwQCMAAwEQYJYIZIAYb4QgEBBAQDAgZAMDMGCWCGSAGG\n+EIBDQQmFiRPcGVuU1NMIEdlbmVyYXRlZCBTZXJ2ZXIgQ2VydGlmaWNhdGUwHQYD\nVR0OBBYEFEEFnOmrROOjNNQrLRoXPXsJLPkeMIGdBgNVHSMEgZUwgZKAFOlVp+B1\nWShwNQuSAVuQAvCe/teZoXakdDByMQswCQYDVQQGEwJBVTEMMAoGA1UECAwDTlNX\nMQ8wDQYDVQQHDAZTeWRuZXkxDTALBgNVBAoMBG15Y2ExEjAQBgNVBAsMCW15Y2Eg\ncm9vdDEhMB8GA1UEAwwYbXkgY2VydGlmaWNhdGUgYXV0aG9yaXR5ggIQADAOBgNV\nHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDQYJKoZIhvcNAQELBQAD\nggIBABmHgp3VwBiQ+viw9exwx8tgTkPhCu36qkDIlSnaPlKTCmlcc+Cs399ra4as\nyVuLPiTGQZ1HsNtmZ3DxaRWbtfdxKty13mZi1+x1UKr/MrGiUTtpYsptJSUYppWa\n7leA1nO/5Kz8i2WCFuk+K3HNVRdjIYmhB3pG3IEXukmaZHSVJ5fCi1ED4l9gzkPv\nTS4olPOU37RPsTgH2ibQUxqhSt0wfu/X6dgqYf3JYtEl4Ddw1XcQeKDI/08D+XP/\nuzNBciMtcAxmTDM+daTBZ8KZnHnPDeuPCj0yLxMi4/HlzuCUXmnO7TAabVyZDvoA\nTpPhkIxj4BuYjCIX9Czd+1fqIu+22tovWg54o+2vuMKyeRpbw3lwTfX3mBuzPaoc\nFu2wFSQEsSQVrpnyD3wMhvF9X2S9YlrzuwQZRJkuYpO83VMHtWIasn6q1al57V0+\nx9TQpCkT6hqv26VxyDhUumAlBkoqEEVXfk/zSa63JHck9LLEIuVt6se6muNj7pLF\nSM+JQe41CQsnxNABs4FqOxp+RhYhOPIKlpgdhBtob89y6OPlR9Sa/4Wchf7FbDPW\nj1ZCzwScOlYDVamud+/wSOsvZhjnkYv8YW13z6GQPqmyMcu0QstIBHkModoZIKOG\nK3FXY7kDCcrj6luxalsD7GbQ7gCpDIdlRe9JPfOBo342mFKR\n-----END CERTIFICATE-----\n"
	keyPem := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA3cUUcNzd9YjLG6QwNGuIRti1UDjUR6JCLB/tEt0PGH43ZUBr\nNRtvFQnTqcoCyVfb7pmYmsyVL6KidfUG8+E/Leye5kjWai9dy8fplXnJQZnIkOhy\nM5Mceo3A9IM4BBfzu0uNBsJrxHD/8Hjzwag9tpmVZgO9lDQzDmtqyK7QH44jo8vz\n5WdzaNd0+Vze2I2ghVdFnXhM16LttQMDc7XLn3eJ+3a1D9NnGVjt/KVlaFNXixSb\n/3KT/Ohv53LOujv3knLy9OGV48A3/+1kYtYXMbtryj0zHfoE+dxJnS/4R8S8/bhl\n98gMVG3KwF7C7b7fchWkhKVtoeCPqey6/wYUzQIDAQABAoIBACPJuhK8kdUdzikX\nxe+vqr5EGn5nrVoiBSu5uzhgFB+PvsDINITNeI+cllvADdMQKp3Gi6nveePGCxGe\nCREyOE/g74OaHX/lRO2txTQqAyBjAMrhuAw6oU3lsk3DHzcJ5ntDJe8BUQLSeXsF\nCdEmpU7iWgmscNuJ0PNywjjAfTWaHXNgXbcragVT/El53/fAnO36aDmd5SP4BiiQ\n984Hig9Z+B9AuqYzKour8o96+IC8eD6EzSVbyvE7WnUZiVV2Opf4mJ8qUEw1NQlg\nGScrcF5RSCJTmB1lt9/mLE1PFS2SZpt2u3iCyKPAqWLa3oAzWMqD9X45+UV2UFlV\nnrfkrsECgYEA/rEE64qKiR5dgjvZVJls6dVu2WYy+EXCSqY2mYFbzHP+rw/xs7oZ\nk39/c0QghZJXDzzxXFUgKa5oeKrkYefPBWFquUfZx/OltbWfjdk8L/z8kfpYJetB\nySELnZiq9mb0JcDPGT5TJVR/udTlCtz89VPeYVt7dOypsAF0uvSrrUcCgYEA3ujC\nvvlughdm7oqhIgaRsIZKQedXLQVb8B1X1HnrbDgnuvBXEKioxIZT6Aw73scl5IFU\n7VBA+tasm9MdwtM18wJ62XCKuN3EgAA0/XpiuageWxSMfwm4Gy2t6FnV5CM+3in/\nmEPDG4NiUqyhk8eDuuuPLWtnXpRN+HQKM5xHp0sCgYEA3dZb/bkXP6WGNxhgDRLx\nzZ6MxakBvkQsng62QfBtf+CMtfjCQxRWkKWd4k01soIreGdRp2Wx9PwnnOrkr+5T\n4FDgv2843rN245XF2qybgwTtDU0rmCOYklJJJsTCLIqyH2wYNtmVXE+ETN2FfnfL\nkPezG8Ot/cLhbh9miCzyl6MCgYEAgYU9oznLvEtcw75JYjvu62McQq7pOH+krCBg\nqFUvNfJrI3QDIurdJVPn7S0unIOawOtlLX80Qov6P5Cr+kg/ULRgLXf3IvO4+acl\nIyO5uaa1/LYz7Jz5HNGt+xQ39BeGsBA3M4IsHBB7UQ591CBZqoK07u85YPtLUtIa\nG2LzP4ECgYBHOPg3ndFMe5EBql/92nSH+RILE6ADUCa+oQUOKa5p/cdWMt6ClT0m\n6cMOJN8lMmtVzwRG/aLPhN2L/vCbtBFDBDIm8PM5gg0340uFv5Mo4p1Sf8iRZG4B\nmzl86a1/OBk4MrtJqoqKrR9yg5/BXlvwuXBJRHaLjGERxhzyhk/WaQ==\n-----END RSA PRIVATE KEY-----\n"
	interCAPem := "-----BEGIN CERTIFICATE-----\nMIIFzDCCA7SgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwcjELMAkGA1UEBhMCQVUx\nDDAKBgNVBAgMA05TVzEPMA0GA1UEBwwGU3lkbmV5MQ0wCwYDVQQKDARteWNhMRIw\nEAYDVQQLDAlteWNhIHJvb3QxITAfBgNVBAMMGG15IGNlcnRpZmljYXRlIGF1dGhv\ncml0eTAgFw0yMDA5MTcyMTE3MjdaGA8yMjM5MDkzMDIxMTcyN1owdjELMAkGA1UE\nBhMCQVUxDDAKBgNVBAgMA05TVzENMAsGA1UECgwEbXljYTEaMBgGA1UECwwRbXlj\nYSBpbnRlcm1lZGlhdGUxLjAsBgNVBAMMJW15IGNlcnRpZmljYXRlIGF1dGhvcml0\neSBpbnRlcm1lZGlhdGUwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCm\n1qRLG0twxwbxBdt/nDeUr0Ia8LtLDtvPjdVVUDTpCp4gnlkEEzHu8JXPPsIami2C\nt5vo35JW1AI223FD5eef54wZG2rXlJbzwlB+yyE5+/V/6WSKe42rePvZDCD+Ym/Q\nyYeqzObViXGnmIvta2aEYZzLeTJPzppvQws/bM+d5IhRa43JuJOVYmjPdp1cjaOm\ntmW3zQSj/00a3i/97SHoyqaJX+y2bPQIJ+yScdBSn9W+Ke3o7/WnuP0HO/ST1fZM\nyzorGbso6aGnTswFbOdWMDUpauE97SL1M6ztoaI4a0HHD8Z8dPhtAmXWbs/5hmQr\njZqBj5W4oUik9iIjUhC2l1aYUf934Om62JjMn9if/mIIA5UTorddj/wKtIsd0n4X\nq5nhJ+X4yVXi3YjqW8iegenaq6UGuvNsm6m/JRAf+5n3FuspHH4WrCgAaIrYg0ZY\nDDu5ro6zHxTcHF6j01CXlJTDEJlStoZ6N9cIKVT94pUPM+EZBq3DGlhBDKipZWk7\n+sEu7sZoQ51WoV4haMY+4Wd7ea8o4sE50eoW+DN2o9lIPHMyxY5uFD7CluUt/b37\ntCcOYAV86JWBN5htTPYAH3wXsDBU/KFSJPLRPF96cuHL6Dq++Gvlqw0rKDKQ/gKh\nDma8lZ9SjVTskqk3l5wzHyNjy7nYFSIRItGIhVbp0wIDAQABo2YwZDAdBgNVHQ4E\nFgQU6VWn4HVZKHA1C5IBW5AC8J7+15kwHwYDVR0jBBgwFoAUD+ANepMk0O9Poxxx\nMpCnxxVyHNMwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAYYwDQYJ\nKoZIhvcNAQELBQADggIBAAQKj6FLBiy23kqHB7iUrl3dSXjJEsPm03zApRhWhr3e\nuxGVYO/YM6RlcJlc7RiKrQAO7XMuOfGbV/TedKPYz+SAeoHCdAVmT21o9HqgwRJ2\nkbJulqIF7oRmmqFOUDIUNg+ZC68QvR9cfuhzcLsEdmfEhXvI5j6CvrhOUN3UHw8A\nO7b4kiymBVT88uXUC0i3bGeEI3h6Fz/RZLbShcvTz2BwcuqoWdInyKi+8mKNfc1O\n+HGBMjnPahNAiovaEuUGErloETdjhmSOkbPBG8h9KpkndCwclEhsBN1+skKiDzKa\nMk53cXXKjqPvPEG9dfQQu0NEnOeY3ZtyVpMqnbo+G0MtyzkozvAB5WjWlpaWZYV2\nnw/wnyCi57ruYI7UjUp+NvFDiIRlOysLC7K6xia+8m7mP8MaFJibQh0tA2UDmdXs\nwy/Z87c6KUCyDB8Hl//rLWbWg6JpHTcH+81yDkVeq2TvJkB6P8jThv51Pz1z4b6U\ndHWAMK5kLmHv+P6sw0JkE5fwszoFOaqSxABq02Pkt5+Hv2EvwxpJZvySkdp7s+Xn\nGUwXhduMscVL/Yd62ES5dYSQ+vbmZIEK3PIttcIyleif6DLFZijJnywf5etYxvrK\nY9wgX6D9PwShl32sf3nzHXh3npLdbio3XwJQUcO6c/lm49rKD7L9L5RM6FNShl8R\n-----END CERTIFICATE-----\n"
	rootCAPem := "-----BEGIN CERTIFICATE-----\nMIIFzDCCA7SgAwIBAgIJAKdYQFPloO6RMA0GCSqGSIb3DQEBCwUAMHIxCzAJBgNV\nBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTENMAsGA1UECgwE\nbXljYTESMBAGA1UECwwJbXljYSByb290MSEwHwYDVQQDDBhteSBjZXJ0aWZpY2F0\nZSBhdXRob3JpdHkwIBcNMjAwOTE3MjA1MTM1WhgPMjI5NDA3MDMyMDUxMzVaMHIx\nCzAJBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5leTENMAsG\nA1UECgwEbXljYTESMBAGA1UECwwJbXljYSByb290MSEwHwYDVQQDDBhteSBjZXJ0\naWZpY2F0ZSBhdXRob3JpdHkwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoIC\nAQDgvRdI24Rv9XBnirlB1LwS32MYyVM2mksTF52E0qrg1OKcMs1D2737BrgaUD6C\nB1I2lMAKR25Q3+x9fSutyww8KZ7yQkFcX2lhwsyYll0j1rvkjek0M1K4787ZFrXS\ncRihE6BSvP5886O+v7a30HxtKbI9oFHdbzgpLpTzvVAn53tokRgAJNtQZWpyJ5Qq\nIG7c96dG9zsXE5+tYT0E0p3ec1z/Ucdx6SKOFjCR8bVLX+Y97mxypOMaPEhGJ4D3\nBlxlCvwDo5sF46e/ntie3Fqghk3jRZTUXedB0IjN8iJCKODPMO1j1cESqVg21xGZ\nyZxIn/ra1iqx9VDCP8egfUOmmMF8flGV08qOGDLGEc/dpVe/yHvG3lmld3MBsW+3\nu6O2l7GIKdLHKibe3uGHhmuPbHq2vlc6IIlRtpsZtK3IXt+bpvlKdI3rxbl4MbT7\n8Z09IUpTsT5jDPEVRnX0zV78Gs4TyKqJKxJJaINx9n0AuXJ8b3jmth/Bb6OkoPgv\nsbFS2QER2Yp8whE1W2PMwtJ06u20YX0RSwuKD+CsnTVmtQwWLBXescCNRH372HwS\nLHO8dvyFWfekLaB2LfciJWYBd8thO5Y4O65FnKLGDvEUh6Ew2OOnhOpy4flWAng6\n39r5uuDQqmWrPFjDNR5HvQjQu1Bv0j81cFY4qZqSIskR9wIDAQABo2MwYTAdBgNV\nHQ4EFgQUD+ANepMk0O9PoxxxMpCnxxVyHNMwHwYDVR0jBBgwFoAUD+ANepMk0O9P\noxxxMpCnxxVyHNMwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAYYwDQYJ\nKoZIhvcNAQELBQADggIBAFdMzocjv6RojMXft1TnwYKb8H0ce6qcsBHZmd/M7IXf\nhyedRkcm7RuN7ayNjFA+44pwAr4jMMNklBQDpGD5yYt1jsltiYoYX5bwZdn2I/rr\nwNQ/FfNSp8rJWqtBhaEt0VI+snHuy0Gdx1eQGf4bJNzvsDLjjJuQ32VUjaCOzsd1\n7d8jR/yjR3Sq20oFEu3HqFSC9OCH2QORTqf6i2IkaUeJbkVTa8+uVceDDbRs3CwY\nVgk/4WcOzcrz0F2BJPpFQ4knrSuHgUbElPHPVuZcn3XZ0n1KBXZdNVCIyLVRowdr\nI+gNEgWE3670Osx55QWg7depP7hU30nQlC1cm2ej2MxM48ddbAL4Zqs8/W1gm+Xb\nDkTsfh81QZQaw6qFVGHJNRIyfMT68ekFB8AgqntulIFR2RJTr/3QJBMhGHKQkmcT\nsa0z0ZrmS/ieurRUjaCsud10Y5VbY5Y8ll5kPsuRWuyijftjcPFqHBzLSSdLacO9\nlVIGkTA3ARCGgym3v5+ZZJ4DeLOJRz9c9OCIASlCkNFFEm1aJ8oagynh2tYqe5TK\nCva1MX8QW5OjHbrm1xvQ8uZOSj55yuBQWKH47GF4QxiojzKikLv4Cpv2Tk5SR9qv\nq3C4t8B26KurNb4z99eo5XhW5XXvQdKZTQC9BqZDN7xhQlwm5lbRSuhZMBJJaQOS\n-----END CERTIFICATE-----\n"

	tlsConfig := mockXTTlsConfig(certPem, keyPem, interCAPem, rootCAPem)
	verified, err := checkCertChain(tlsConfig.Certificates[0])
	logCertStats(verified)
	if err != nil {
		t.Errorf("certificate chain with 1 TLS cert, 1 intermediate, 1 root cert not validated, cause: %s", err)
	} else {
		t.Logf("normal. certificate chain with 1 TLS cert, 1 intermediate, 1 root cert validated, length: %d", len(verified))
	}
}

func mockXTTlsConfig(certPem string, keyPem string, intercAPem string, rootcAPem string) *tls.Config {
	r := mockRuntime()
	r.Connection.Downstream.Cert = certPem + intercAPem + rootcAPem
	r.Connection.Downstream.Key = keyPem
	r.Connection.Downstream.Mode = "TLS"
	r.Connection.Downstream.Port = 8443

	tlsConfig := r.tlsConfig()
	return tlsConfig
}
