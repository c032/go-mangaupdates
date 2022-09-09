package mangaupdates_test

import (
	"testing"

	"github.com/c032/go-mangaupdates"
)

func TestClient_Time(t *testing.T) {
	muc := mangaupdates.Client{}

	tr, err := muc.Time()
	if err != nil {
		t.Fatal(err)
	}

	if tr == nil {
		t.Fatalf("muc.Time() = nil; want non-nil")
	}

	if got, notWant := tr.Timestamp, int64(0); got == notWant {
		t.Errorf("muc.Time().Timestamp = %#v; want != %#v", got, notWant)
	}

	if got, notWant := tr.AsRFC3339, ""; got == notWant {
		t.Errorf("muc.Time().AsRFC3339 = %#v; want != %#v", got, notWant)
	}

	if got, notWant := tr.AsString, ""; got == notWant {
		t.Errorf("muc.Time().AsString = %#v; want != %#v", got, notWant)
	}
}
