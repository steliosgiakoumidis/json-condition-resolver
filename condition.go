package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

var plot_data []node

type node struct {
	level       int
	item_number int
	value       string
	connects_to int
}

func main() {

	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	operator, operand, err := get_object(string(bytes), "{", "}")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	operator = extract_operator(operator)
	operand = extract_operand(operand)

	resp, err := process_condition(operator, operand, 0, 0)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	fmt.Printf("%t\n\n", resp)

	plot_tree(plot_data)
}

func extract_operator(input string) (operator string) {
	if strings.Contains(input, "!") {
		return input[1:2]
	} else {
		return input[1:3]
	}
}

func extract_operand(input string) (operand string) {
	return input[1 : len(input)-1]
}

func plot_tree(plot_data []node) {
	sort.Slice(plot_data, func(i, j int) bool {
		if plot_data[i].level < plot_data[j].level {
			return true
		}
		if plot_data[i].level > plot_data[j].level {
			return false
		}

		return plot_data[i].connects_to < plot_data[j].connects_to
	})

	current_level := 0
	for _, v := range plot_data {
		if v.level != current_level {
			fmt.Printf("\n")
			current_level = v.level
		}
		fmt.Printf("%+v \t", v.value)

	}
}

func get_last_item_of_level(level_number int) (last_object_value int) {
	max_item_value := 0
	for _, val := range plot_data {
		if val.level == level_number {
			if val.item_number > max_item_value {
				max_item_value = val.item_number
			}
		}
	}

	return max_item_value
}

func all(slice []bool) (result bool) {
	for _, v := range slice {
		if !v {
			return false
		}
	}

	return true
}

func any(slice []bool) (result bool) {
	for _, v := range slice {
		if v {
			return true
		}
	}

	return false
}

func get_object(str string, start string, end string) (operator string, operant string, err error) {
	start_index := strings.Index(str, start)
	if start_index == -1 {
		return "", "", errors.New("object start cannot be found")
	}
	start_index += len(start)
	end_index := strings.LastIndex(str, end)
	if end_index == -1 {
		return "", "", errors.New("object end cannot be found")
	}
	end_index = start_index + end_index - 1

	objectParsed := strings.SplitN(str[start_index:end_index], ":", 2)

	return objectParsed[0], objectParsed[1], nil
}

func process_condition(operator string, operand string, tree_level int, connects_to int) (local_result bool, err error) {

	switch operator {
	case "==":
		item_number := get_last_item_of_level(tree_level) + 1
		operands := strings.Split(operand, ",")

		if len(operands) != 2 {
			return false, errors.New("equality condition is not properly structured")
		}

		plot_data = append(plot_data, node{level: tree_level, value: operator, item_number: item_number, connects_to: connects_to})

		connects_to = item_number
		next_level_item_number := get_last_item_of_level(tree_level+1) + 1
		for _, element := range operands {
			plot_data = append(plot_data, node{level: tree_level + 1, value: element, item_number: next_level_item_number, connects_to: connects_to})
			next_level_item_number++
		}

		if operands[0] == operands[1] {
			return true, nil
		}

		return false, nil

	case "&&", "||":
		item_number := get_last_item_of_level(tree_level) + 1
		plot_data = append(plot_data, node{level: tree_level, value: operator, item_number: item_number, connects_to: connects_to})

		operands := strings.Split(operand, "},{")

		var new_objects []string
		for i, val := range operands {
			if i == 0 {
				new_objects = append(new_objects, val+"}")
				continue
			}
			if i == len(operands)-1 {
				new_objects = append(new_objects, "{"+val)
				continue
			}
			new_objects = append(new_objects, "{"+val+"}")
		}

		var loop_results []bool
		connects_to = item_number
		for _, el := range new_objects {
			local_operator, local_operand, err := get_object(el, "{", "}")
			if err != nil {
				return false, err
			}
			local_operator = extract_operator(local_operator)
			local_operand = extract_operand(local_operand)

			local_result, err := process_condition(local_operator, local_operand, tree_level+1, connects_to)
			if err != nil {
				return false, err
			}

			loop_results = append(loop_results, local_result)
		}

		if operator == "&&" && all(loop_results) {
			return true, nil
		}

		if operator == "||" && any(loop_results) {
			return true, nil
		}

		return false, nil

	case "!":
		item_number := get_last_item_of_level(tree_level) + 1
		plot_data = append(plot_data, node{level: tree_level, value: operator, item_number: item_number, connects_to: connects_to})

		local_operator, local_operand, err := get_object(operand, "{", "}")
		if err != nil {
			return false, err
		}

		local_operator = extract_operator(local_operator)
		local_operand = extract_operand(local_operand)

		connects_to = item_number
		resp, err := process_condition(local_operator, local_operand, tree_level+1, connects_to)
		if err != nil {
			return false, err
		}

		return !resp, nil

	default:
		return false, errors.New("no operator match")
	}
}
