package yamlhelper

import (
	"github.com/iPanja/go-yamlx/parser"
)

func Read(filename string, out interface{}) error {
	return parser.ReadYAML(filename, out)
}

func Write(filename string, in interface{}) error {
	return parser.WriteYAML(filename, in)
}
