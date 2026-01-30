package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexValue aceita string, number ou null do JSON
// Usado para campos que o PHP envia como string vazia ou número
type FlexValue string

func (f *FlexValue) UnmarshalJSON(data []byte) error {
	// Tenta string primeiro
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexValue(s)
		return nil
	}

	// Tenta número inteiro
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexValue(fmt.Sprintf("%d", i))
		return nil
	}

	// Tenta número float
	var fl float64
	if err := json.Unmarshal(data, &fl); err == nil {
		*f = FlexValue(fmt.Sprintf("%.1f", fl))
		return nil
	}

	// Tenta null
	if string(data) == "null" {
		*f = ""
		return nil
	}

	// Fallback: usa como string literal
	*f = FlexValue(string(data))
	return nil
}

func (f FlexValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexValue) String() string {
	return string(f)
}

// Float converte FlexValue para float64 (retorna 0 se não for numérico)
func (f FlexValue) Float() float64 {
	s := string(f)
	if s == "" || s == "-" {
		return 0
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// FlexInt aceita int ou null do JSON
type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = 0
		return nil
	}

	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		// Tenta string com número
		var s string
		if err2 := json.Unmarshal(data, &s); err2 == nil {
			fmt.Sscanf(s, "%d", &i)
		}
	}
	*f = FlexInt(i)
	return nil
}

func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(f))
}

func (f FlexInt) Int() int {
	return int(f)
}

// FlexBool aceita bool, int (0/1), ou string ("true"/"false") do JSON
type FlexBool bool

func (f *FlexBool) UnmarshalJSON(data []byte) error {
	// Tenta bool primeiro
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*f = FlexBool(b)
		return nil
	}

	// Tenta int (0 = false, != 0 = true)
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexBool(i != 0)
		return nil
	}

	// Tenta string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexBool(s == "true" || s == "1")
		return nil
	}

	// Null = false
	if string(data) == "null" {
		*f = false
		return nil
	}

	return nil
}

func (f FlexBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(f))
}

func (f FlexBool) Bool() bool {
	return bool(f)
}
