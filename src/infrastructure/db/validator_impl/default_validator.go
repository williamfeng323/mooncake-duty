package validatorimpl

import (
	"fmt"
	"reflect"
	"strconv"
)

//DefaultValidator will validate below tag:
//    required: boolean <-- The field could not be empty.
//    TBD: More default validators
type DefaultValidator struct{}

// Verify returns empty error list if validate successfully.
func (dv DefaultValidator) Verify(model interface{}) []error {

	docElem := reflect.ValueOf(model).Elem()
	fieldType := docElem.Type()
	validationErrors := []error{}

	for fieldIndex := 0; fieldIndex < docElem.NumField(); fieldIndex++ {
		var required bool
		var err error
		var fieldValue reflect.Value
		field := fieldType.Field(fieldIndex)
		fieldTag := field.Tag
		requiredTag := fieldTag.Get("required")
		fieldElem := docElem.Field(fieldIndex)
		fieldName := field.Name
		if fieldElem.Kind() == reflect.Ptr || fieldElem.Kind() == reflect.Interface {
			fieldValue = fieldElem.Elem()
		} else {
			fieldValue = fieldElem
		}

		if len(requiredTag) > 0 {
			required, err = strconv.ParseBool(requiredTag)
			if err != nil {
				panic("Check your required tag - must be boolean")
			}
			if required {
				if err := validateRequired(fieldValue, fieldName); err != nil {
					validationErrors = append(validationErrors, err)
				}
			}
		}
		if len(validationErrors) > 0 {
			return validationErrors
		}
	}
	return nil
}

func validateRequired(fieldValue reflect.Value, fieldName string) error {
	isSet := false
	if !fieldValue.IsValid() {
		isSet = false
	} else if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
		isSet = true
	} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Map {
		isSet = fieldValue.Len() > 0
	} else if fieldValue.Kind() == reflect.Interface {
		isSet = fieldValue.Interface() != nil
	} else {
		va := fieldValue.Interface()
		zeroValue := reflect.Zero(reflect.TypeOf(fieldValue.Interface())).Interface()
		isSet = !reflect.DeepEqual(va, zeroValue)
	}

	if !isSet {
		return fmt.Errorf("Field - %s cannot be null", fieldName)
	}
	return nil
}

// NewDefaultValidator returns an empty NewDefaultValidator instance.
func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}
