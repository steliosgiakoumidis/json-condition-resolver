package main

import (
	"testing"
)

func TestAllFuncReturnsTrueWhenAllElementsAreTrue(t *testing.T) {
	slice := []bool{true, true, true}
	expected_value := true
	actual_value := all(slice)
	if actual_value != expected_value {
		t.Fatalf("Expected value %v and actual value is %v", expected_value, actual_value)
	}
}

func TestAllFuncReturnsFalseWhenOneElementsIsFalse(t *testing.T) {
	slice := []bool{true, false, true}
	expected_value := false
	actual_value := all(slice)
	if actual_value != expected_value {
		t.Fatalf("Expected value %v and actual value is %v", expected_value, actual_value)
	}
}

func TestAnyFuncReturnsTrueWhenOneElementsIsTrue(t *testing.T) {
	slice := []bool{true, false, false}
	expected_value := true
	actual_value := any(slice)
	if actual_value != expected_value {
		t.Fatalf("Expected value %v and actual value is %v", expected_value, actual_value)
	}
}

func TestGetLastItemOfLevelFuncReturnsTheLargestValue(t *testing.T) {
	input := "{\"||\":[{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	_, _ = process_condition(operator, operand, 0, 0)

	expected_value := 4
	actual_value := get_last_item_of_level(2)
	if actual_value != expected_value {
		t.Fatalf("Actual value=%v, Expected value=%v", actual_value, expected_value)
	}
}

func TestGetObjectReturnsOperandAndOperator(t *testing.T) {
	test_object := "{\"==\":[5,4]}"
	operator, operand, err := get_object(test_object, "{", "}")
	expected_operator := "\"==\""
	expected_operand := "[5,4]"

	if operator != expected_operator || operand != expected_operand || err != nil {
		t.Fatalf("Actual operator=%v, Expected operator=%v, Actual operand=%v, Expected operand=%v, Error=%v", operator, expected_operator, operand, expected_operand, err)
	}
}

func TestGetObjectReturnsErrorBecauseStartCharMissing(t *testing.T) {
	test_object := "\"==\":[5,4]}"
	_, _, err := get_object(test_object, "{", "}")
	expected_error := "object start cannot be found"

	if err.Error() != expected_error {
		t.Fatalf("Expected error is %v and actual error is %v", expected_error, err.Error())
	}
}

func TestGetObjectReturnsErrorBecauseEndCharMissing(t *testing.T) {
	test_object := "{\"==\":[5,4]"
	_, _, err := get_object(test_object, "{", "}")
	expected_error := "object end cannot be found"

	if err.Error() != expected_error {
		t.Fatalf("Expected error is %v and actual error is %v", expected_error, err.Error())
	}
}

func TestProcessConditionSuccessfulWithProperObject(t *testing.T) {
	input := "{\"||\":[{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	result, _ := process_condition(operator, operand, 0, 0)
	expected_result := true

	if result != expected_result {
		t.Fatalf("Expected result is %v and actual result is %v", expected_result, result)
	}
}

func TestProcessConditionSuccessfulWithThreeConditions(t *testing.T) {
	input := "{\"||\":[{\"==\":[\"foo\",\"bar\"]},{\"==\":[\"log\",\"log\"]},{\"!\":[{\"==\":[1,1]}]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	result, _ := process_condition(operator, operand, 0, 0)
	expected_result := true

	if result != expected_result {
		t.Fatalf("Expected result is %v and actual result is %v", expected_result, result)
	}
}

func TestProcessConditionCorrectlySetFalseWithThreeConditions(t *testing.T) {
	input := "{\"&&\":[{\"==\":[\"foo\",\"foo\"]},{\"==\":[\"log\",\"log\"]},{\"!\":[{\"==\":[1,1]}]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	result, _ := process_condition(operator, operand, 0, 0)
	expected_result := false

	if result != expected_result {
		t.Fatalf("Expected result is %v and actual result is %v", expected_result, result)
	}
}

func TestProcessConditionFailsOnMalformedObject(t *testing.T) {
	input := "{\"||\":[{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1,1]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	_, err := process_condition(operator, operand, 0, 0)
	expected_error := "equality condition is not properly structured"

	if err.Error() != expected_error {
		t.Fatalf("Expected error is %v and actual error is %v", expected_error, err.Error())
	}
}

func TestProcessConditionFailsNoMatchingOperator(t *testing.T) {
	input := "{\"?\":[{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1]}]}"
	operator, operand, _ := get_object(string(input), "{", "}")

	operator = operator[1:3]
	operand = operand[1 : len(operand)-1]
	_, err := process_condition(operator, operand, 0, 0)
	expected_error := "no operator match"

	if err.Error() != expected_error {
		t.Fatalf("Expected error is %v and actual error is %v", expected_error, err.Error())
	}
}

func TestExtractOperand(t *testing.T) {
	input := "[{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1]}]"
	actual_operand := extract_operand(input)
	expected_operand := "{\"==\":[\"foo\",\"bar\"]},{\"==\":[1,1]}"

	if actual_operand != expected_operand {
		t.Fatalf("Expected operand is %v and actual operand is %v", expected_operand, actual_operand)
	}
}

func TestExtractOperator(t *testing.T) {
	input := "\"!\""
	actual_operator := extract_operator(input)
	expected_operator := "!"

	if actual_operator != expected_operator {
		t.Fatalf("Expected operator is: %v and actual operator is %v", actual_operator, expected_operator)
	}
}
