package utils

import (
	"testing"
	"guthub.com/Toront0/lux-server/internal/utils"
)

func TestCamelCaseToSnakeCase(t *testing.T) {

	res := utils.CamelCaseToSnakeCase("firstName")


	if res != "first_name" {
		t.Errorf("CamelCaseToSnakeCase FAILED! Expected %s, got %s", "first_name", res)
	}

}