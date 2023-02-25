package egorm

func DbGetAll[T any](input *[]T) error {

	if err := InitDB(); err != nil {
		return err
	}

	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}
	result := Db.Find(input)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DbCreate[T any](input *T) error {
	if err := InitDB(); err != nil {
		return err
	}
	err := autoMigrate(input)
	if err != nil {
		return err
	}
	result := Db.Create(input)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DbGet[T any](input []T, where map[string]string) error {
	if err := InitDB(); err != nil {
		return err
	}
	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}
	return nil
}

func DbFirst[T any](input *T, where map[string]string) error {
	if err := InitDB(); err != nil {
		return err
	}
	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}
	return nil
}

func DbFind[T any](input *T, id int) error {
	if err := InitDB(); err != nil {
		return err
	}
	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}
	return nil
}
