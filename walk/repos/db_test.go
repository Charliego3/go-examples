package repos

import "testing"

func TestFetchUser(t *testing.T) {
	user, err := FetchUser()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("User: %+v", user)
}

func TestFetchAreas(t *testing.T) {
	areas, err := FetchAreas()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Areas: %#v", areas)

	// Fetch tables
	tables, err := FetchTables(areas[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Tables: %#v", tables)
}
