package validation

import "errors"

func IsValidField(field string, value any) error {
	validFunc := map[string]func(value any) error{
		"description": isValidDescription,
		"rate":        isValidRate,
		"name":        isValidName,
	}
	f, ok := validFunc[field]
	if !ok {
		return errors.New("invalid field")
	}
	return f(value)
}

func isValidDescription(value any) error {
	description, ok := value.(string)
	if !ok {
		return errors.New("description is not string")
	}
	runeLen := len([]rune(description))
	if runeLen <= 1000 {
		return nil
	}
	return ErrInvalidDescription
}

func isValidRate(value any) error {
	rate, ok := value.(int)
	if !ok {
		return errors.New("rate is not int")
	}
	if rate < 0 || rate > 10 {
		return ErrInvalidRate
	}
	return nil
}
func isValidName(value any) error {
	name, ok := value.(string)
	if !ok {
		return errors.New("name is not string")
	}
	runeLen := len([]rune(name))
	if runeLen > 150 || runeLen < 1 {
		return ErrInvalidName
	}
	return nil
}
