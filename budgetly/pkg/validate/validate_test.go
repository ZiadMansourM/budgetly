package validate

import (
	"testing"
)

func TestMainValidator(t *testing.T) {
	// Sample data to validate
	data := struct {
		FirstName string
		LastName  string
		Email     string
	}{
		FirstName: "",         // Should trigger the "required" rule
		LastName:  "Mansour",  // Should pass validation
		Email:     "invalid@", // Should fail email validation
	}

	// Define the validation rules for each field
	validationFields := ValidationFields{
		"FirstName": Rules(
			Required,
			Min(3),
			ErrorMessage("First name is required and must have at least 3 characters"),
		),
		"LastName": Rules(
			Required,
			Min(3),
			ErrorMessage("Last name is required and must have at least 3 characters"),
		),
		"Email": Rules(
			Required,
			Email,
			ErrorMessage("Please enter a valid email address"),
		),
	}

	// Call the Validate function to validate the data
	errors := Validate(data, validationFields)

	// Check if there are any errors
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d: %v", len(errors), errors)
	}

	// Verify the error messages for FirstName and Email
	if errors["FirstName"] != "First name is required and must have at least 3 characters" {
		t.Errorf("Expected custom error message for FirstName, got %s", errors["FirstName"])
	}

	if errors["Email"] != "Please enter a valid email address" {
		t.Errorf("Expected custom error message for Email, got %s", errors["Email"])
	}
}

func TestRequired(t *testing.T) {
	rule := Required()
	rule.FieldName = "username"
	rule.FieldValue = ""

	if rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to be required", rule.FieldName)
	}

	rule.FieldValue = "Ziad"
	if !rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to pass validation when not empty", rule.FieldName)
	}
}

func TestEmail(t *testing.T) {
	rule := Email()
	rule.FieldName = "email"
	rule.FieldValue = "invalid-email"

	if rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to fail validation for invalid email", rule.FieldValue)
	}

	rule.FieldValue = "test@example.com"
	if !rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to pass validation for valid email", rule.FieldValue)
	}
}

func TestMax(t *testing.T) {
	rule := Max(5)()
	rule.FieldName = "username"
	rule.FieldValue = "abcdef"

	if rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to fail max length validation", rule.FieldValue)
	}

	rule.FieldValue = "abc"
	if !rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to pass max length validation", rule.FieldValue)
	}
}

func TestMin(t *testing.T) {
	rule := Min(3)()
	rule.FieldName = "username"
	rule.FieldValue = "ab"

	if rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to fail min length validation", rule.FieldValue)
	}

	rule.FieldValue = "abc"
	if !rule.ValidationFunc(rule) {
		t.Errorf("Expected %s to pass min length validation", rule.FieldValue)
	}
}
