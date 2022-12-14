package mangaupdates_test

import (
	"testing"

	"github.com/c032/go-mangaupdates"
)

func TestClient_Time(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

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

func TestClient_SeriesSearch(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	muc := mangaupdates.Client{}

	req := mangaupdates.SeriesSearchRequest{
		Search: "new game",
		SType:  mangaupdates.STypeTitle,
	}

	resp, err := muc.SeriesSearch(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		t.Fatalf("muc.SeriesSearch(req) = nil; want non-nil")
	}

	if got := len(resp.Results); got == 0 {
		t.Errorf("len(muc.SeriesSearch(req)) = %#v; want non-zero", got)
	} else {
		for i, result := range resp.Results {
			if got := result.Record; got == nil {
				t.Errorf("muc.SeriesSearch(req)[%d].Record = %#v; want non-zero", i, got)
			} else {
				record := result.Record
				if got := record.SeriesID; got == 0 {
					t.Errorf("muc.SeriesSearch(req)[%d].Record.SeriesID = %#v; want non-zero", i, got)
				}
				if got := record.Title; got == "" {
					t.Errorf("muc.SeriesSearch(req)[%d].Record.Title = %#v; want non-empty string", i, got)
				}
				if got := record.URL; got == "" {
					t.Errorf("muc.SeriesSearch(req)[%d].Record.URL = %#v; want non-empty string", i, got)
				}
			}
		}
	}
}

func TestClient_SeriesByID(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	muc := mangaupdates.Client{}

	seriesID := mangaupdates.SeriesID(24348008396)

	series, err := muc.SeriesByID(seriesID)
	if err != nil {
		t.Fatal(err)
	}

	if series == nil {
		t.Fatalf("muc.SeriesByID(%d) = nil; want non-nil", seriesID)
	}

	if got := series.SeriesID; got == 0 {
		t.Errorf("muc.SeriesByID(%d).SeriesID = %#v; want non-zero", seriesID, got)
	}
	if got, want := series.Title, "Youjo Senki (Novel)"; got != want {
		t.Errorf("muc.SeriesByID(%d).Title = %#v; want %#v", seriesID, got, want)
	}
	if got, want := series.URL, "https://www.mangaupdates.com/series/b6o6bik/youjo-senki-novel"; got != want {
		t.Errorf("muc.SeriesSearch(%d).URL = %#v; want %#v", seriesID, got, want)
	}
	if got, want := series.Type, "Novel"; got != want {
		t.Errorf("muc.SeriesSearch(%d).Type = %#v; want %#v", seriesID, got, want)
	}
}
