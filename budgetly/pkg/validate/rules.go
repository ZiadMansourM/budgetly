package validate

import (
	"fmt"
	"regexp"
)

// ErrorMessage allows users to specify custom error messages in the validation chain.
func ErrorMessage(message string) ValidationRuleFunc {
	return func() ValidationRule {
		return ValidationRule{
			Name: "custom_message",
			ErrorMessageFunc: func(rule ValidationRule) string {
				return message
			},
			ValidationFunc: func(rule ValidationRule) bool {
				// This rule doesn't perform validation but modifies the error message
				return true
			},
		}
	}
}

// Required rule ensures that the field is not empty.
func Required() ValidationRule {
	return ValidationRule{
		Name: "required",
		ErrorMessageFunc: func(rule ValidationRule) string {
			return fmt.Sprintf("%s is required", rule.FieldName)
		},
		ValidationFunc: func(rule ValidationRule) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return str != ""
		},
	}
}

// Email rule validates whether the field is a valid email address.
func Email() ValidationRule {
	// emailRegex defines a basic regex for email validation
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return ValidationRule{
		Name: "email",
		ErrorMessageFunc: func(rule ValidationRule) string {
			return fmt.Sprintf("%s is not a valid email", rule.FieldName)
		},
		ValidationFunc: func(rule ValidationRule) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(str)
		},
	}
}

// Max rule validates that the field has a maximum length of `max`.
func Max(max int) ValidationRuleFunc {
	return func() ValidationRule {
		return ValidationRule{
			Name:      "max",
			RuleValue: max,
			ErrorMessageFunc: func(rule ValidationRule) string {
				return fmt.Sprintf("%s must be less than %d characters", rule.FieldName, rule.RuleValue)
			},
			ValidationFunc: func(rule ValidationRule) bool {
				str, ok := rule.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) <= rule.RuleValue.(int)
			},
		}
	}
}

// Min rule validates that the field has a minimum length of `min`.
func Min(min int) ValidationRuleFunc {
	return func() ValidationRule {
		return ValidationRule{
			Name:      "min",
			RuleValue: min,
			ErrorMessageFunc: func(rule ValidationRule) string {
				return fmt.Sprintf("%s must be at least %d characters", rule.FieldName, rule.RuleValue)
			},
			ValidationFunc: func(rule ValidationRule) bool {
				str, ok := rule.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) >= rule.RuleValue.(int)
			},
		}
	}
}
