package egorm

import (
	"fmt"
	"reflect"

	"github.com/martenwallewein/easy-going/pkg/eslices"
)

var alreadyMigratedTypes []string = make([]string, 0)

func getInterfaceTypeAsString[T any](input *T) string {
	return reflect.TypeOf((*T)(nil)).Elem().Name()
}

func autoMigrate[T any](input *T) error {
	typeName := getInterfaceTypeAsString(input)

	if eslices.IndexOf(typeName, alreadyMigratedTypes) >= 0 {
		return nil
	}
	err := Db.AutoMigrate(input)
	if err != nil {
		return fmt.Errorf("egorm: Failed to perform automigration for %s: %s", typeName, err)
	}

	alreadyMigratedTypes = eslices.AppendToSliceIfMissing(alreadyMigratedTypes, typeName)
	return nil
}
