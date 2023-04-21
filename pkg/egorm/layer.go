package egorm

type Where map[string]interface{}

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

func DbSave[T any](input *T) error {
	if err := InitDB(); err != nil {
		return err
	}
	err := autoMigrate(input)
	if err != nil {
		return err
	}
	result := Db.Save(input)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DbGet[T any](input *[]T, where map[string]interface{}) error {
	if err := InitDB(); err != nil {
		return err
	}
	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}

	result := Db.Where(where).Find(input)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DbFirst[T any](input *T, where map[string]interface{}) error {
	if err := InitDB(); err != nil {
		return err
	}
	var tmp T
	err := autoMigrate(&tmp)
	if err != nil {
		return err
	}

	result := Db.Where(where).First(input)
	if result.Error != nil {
		return result.Error
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

	result := Db.Find(input, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
