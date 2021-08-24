package lib

import (
	"fmt"
	"strings"
)

var (
	inputConverterMap  = make(map[string]InputConverter)
	outputConverterMap = make(map[string]OutputConverter)
)

func ListInputConverter() {
	fmt.Println("All available input formats:")
	for name := range inputConverterMap {
		fmt.Printf("  - %s\n", name)
	}
}

func RegisterInputConverter(name string, c InputConverter) error {
	name = strings.TrimSpace(name)
	if _, ok := inputConverterMap[name]; ok {
		return ErrDuplicatedConverter
	}
	inputConverterMap[name] = c
	return nil
}

func ListOutputConverter() {
	fmt.Println("All available output formats:")
	for name := range outputConverterMap {
		fmt.Printf("  - %s\n", name)
	}
}

func RegisterOutputConverter(name string, c OutputConverter) error {
	name = strings.TrimSpace(name)
	if _, ok := outputConverterMap[name]; ok {
		return ErrDuplicatedConverter
	}
	outputConverterMap[name] = c
	return nil
}
