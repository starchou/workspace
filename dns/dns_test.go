package dns

import "testing"

func TestGetPublicIp(t *testing.T) {

	ip, err := getPublicIp()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(ip)
	}
}
