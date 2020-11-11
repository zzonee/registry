package direct

import "testing"

func Test_direct(t *testing.T) {
	dis := NewDiscovery()
	ch, err := dis.Discover("127.0.0.1:8000,127.0.0.1:9000")
	if err != nil {
		t.Fail()
		return
	}
	inss := <-ch
	t.Log(inss)
}