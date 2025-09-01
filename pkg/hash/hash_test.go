package hash

import "testing"

func TestHash(t *testing.T) {
	ans := GetHash("images/221118-N-ON904-1070_52691439736.jpg")
	if ans != 1314607733 {
		t.Errorf(`Hash("images/221118-N-ON904-1070_52691439736.jpg") = %d; want 1314607733`, ans)
	}
}
