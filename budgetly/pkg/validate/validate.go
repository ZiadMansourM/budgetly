package validate

import (
	"reflect"
	"strings"
)

// ValidationRuleFunc defines a function that returns a ValidationRule.
// It is used when a validation rule requires a parameter (e.g., Min, Max).
type ValidationRuleFunc func() ValidationRule

// ValidationRule defines the structure of a validation rule, containing its name, value,
// field value, field name, error message function, and the validation function.
type ValidationRule struct {
	Name             string
	RuleValue        any
	FieldValue       any
	FieldName        any
	ErrorMessageFunc func(ValidationRule) string
	ValidationFunc   func(ValidationRule) bool
}

// ValidationFields maps field names to their associated validation rules.
type ValidationFields map[string][]ValidationRule

// Rules is a variadic function that takes ValidationRuleFuncs and returns a slice of ValidationRules.
// It is used to combine multiple rules for a field.
func Rules(rules ...ValidationRuleFunc) []ValidationRule {
	ruleSets := make([]ValidationRule, len(rules))
	for i, rule := range rules {
		ruleSets[i] = rule()
	}
	return ruleSets
}

// Validator holds the data to be validated and the associated validation rules.
type Validator struct {
	data   any
	fields ValidationFields
}

// New creates a new Validator instance, containing the data and validation rules for fields.
func New(data any, fields ValidationFields) *Validator {
	return &Validator{
		data:   data,
		fields: fields,
	}
}

// Validate applies the validation rules to the fields and returns a map of error messages.
// It uses reflection to extract field values from the data.
func Validate(data any, fields ValidationFields) map[string]string {
	errors := make(map[string]string)

	// Iterate through each field and its associated validation rules
	for fieldName, rules := range fields {
		// Extract the field value from the data using reflection
		fieldValue := getFieldValue(data, fieldName)
		var errorMessage string
		var validationFailed bool

		for _, rule := range rules {
			rule.FieldName = fieldName
			rule.FieldValue = fieldValue

			// Apply the validation rule and use the default error message if it fails
			if !rule.ValidationFunc(rule) {
				validationFailed = true
				errorMessage = rule.ErrorMessageFunc(rule)
			}

			// If it's an ErrorMessage rule, store it but only apply it if a validation rule failed
			if rule.Name == "custom_message" && validationFailed {
				errorMessage = rule.ErrorMessageFunc(rule)
			}
		}

		// If validation failed, store the final error message (custom or default)
		if validationFailed {
			errors[fieldName] = errorMessage
		}
	}

	return errors
}

// Helper function to extract field values from struct using reflection.
// Assumes data is a struct and fields are exported.
func getFieldValue(data any, fieldName string) any {
	dataValue := reflect.ValueOf(data)
	fieldValue := dataValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil
	}
	return fieldValue.Interface()
}

// ValidationError represents validation errors.
type ValidationError struct {
	Errors map[string]string
}

// Error implements the error interface for ValidationError.
// It formats the error message with each field error on a new line.
func (v *ValidationError) Error() string {
	var sb strings.Builder

	sb.WriteString("Validation failed for the following fields:\n")

	for field, err := range v.Errors {
		sb.WriteString(" - ")
		sb.WriteString(field)
		sb.WriteString(": ")
		sb.WriteString(err)
		sb.WriteString("\n")
	}

	return sb.String()
}
