package cmd

import (
	file_util "flag-example/io/file"
	"path/filepath"
	"runtime"
)

// 기본 디렉토리 설정
func DefaultDataDir() string {
	home := file_util.HomeDir()
	// home 디렉토리 받아온 후 이걸로 Default Data Dir 설정하기
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Curie")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Local", "Curie")
		} else {
			return filepath.Join(home, ".curie")
		}
	}

	return ""
}
