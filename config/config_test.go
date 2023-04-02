package config

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var wants [5]string

	//	テストデータ読み込み
	f , err := os.Open("../testdata/config_test/config_test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data := csv.NewReader(f)
	records, err := data.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	for i,record := range records {
		if i == 0 {continue}
		t.Setenv(record[0],record[1])
		wants[i-1] = record[1]
	}

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v",err)
	}
	
	if strings.Compare(got.Env,wants[0]) != 0 {
		t.Errorf("want %s, but %s",wants[0],got.Env)
	}

	if strings.Compare(strconv.Itoa(got.Port),wants[1]) != 0 {
		t.Errorf("want %s, but %d",wants[1],got.Port)
	}


	if strings.Compare(got.RedisHost,wants[2]) != 0 {
		t.Errorf("want %s, but %s",wants[2],got.RedisHost)
	}

	if strings.Compare(strconv.Itoa(got.RedisPort),wants[3]) != 0 {
		t.Errorf("want %s, but %d",wants[3],got.RedisPort)
	}
}