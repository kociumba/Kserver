package luaintegration

import (
	"errors"
	"os"
)

func CreateTempLUA() *os.File {
	tempFile, err := os.CreateTemp("", "kserver.lua")
	if err != nil {
		panic(err)
	}

	if _, err := tempFile.WriteString(Simplify); err != nil {
		panic(err)
	}

	userConfig, err := os.ReadFile("kserver.lua")
	if errors.Is(err, os.ErrNotExist) {
		userConfig = []byte{}
		err = nil
	}
	if err != nil {
		panic(err)
	}
	if _, err := tempFile.Write(append([]byte("\n"), userConfig...)); err != nil {
		panic(err)
	}

	tempFile.Close()
	// defer os.Remove(tempFile.Name())
	return tempFile
}
