package egorm

import (
	"fmt"
	"os"
	"path"
	"testing"

	"gorm.io/gorm"
)

type SampleStruct struct {
	gorm.Model
	Name string
}

func fakeSetup() {

	SetSQLiteConnectOpts(&SQLiteConnectOpts{
		Path: path.Join(os.TempDir(), "egorm_test.sqlite"),
	})
}

func cleanUp() {
	err := os.RemoveAll(path.Join(os.TempDir(), "egorm_test.sqlite"))
	if err != nil {
		fmt.Println(err)
	}
}

func TestSetup(t *testing.T) {
	t.Run("TestSetup", func(t *testing.T) {
		fakeSetup()
		err := InitDB()
		if err != nil {
			t.Error(err)
		}
	})
}

func TestInsert(t *testing.T) {
	t.Run("TestInsert", func(t *testing.T) {
		fakeSetup()
		sample1 := SampleStruct{
			Name: "Sample1",
		}
		err := DbCreate(&sample1)
		if err != nil {
			t.Error(err)
			return
		}

		sample2 := SampleStruct{
			Name: "Sample2",
		}
		err = DbCreate(&sample2)
		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestGetAll(t *testing.T) {
	t.Run("TestGetAll", func(t *testing.T) {
		fakeSetup()

		var samples []SampleStruct
		err := DbGetAll(&samples)
		if err != nil {
			t.Error(err)
			return
		}

		if len(samples) != 2 {
			t.Error(fmt.Errorf("egorm: Expected %d items, got %d", 2, len(samples)))
		}
		cleanUp()
	})
}

/*func TestCleanup(t *testing.T) {

}*/
