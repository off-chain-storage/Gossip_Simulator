package file_util

import (
	"bytes"
	"errors"
	"flag-example/config/params"
	"io"
	"mime/multipart"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 파일의 경로 중 상대 경로로 작성된 문자열을 변환시키기
// ex) ~/aaa/bbb/ccc.c
// -> /home/jinbum/aaa/bbb/ccc.c
func ExpandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
		if home := HomeDir(); home != "" {
			p = home + p[1:]
		}
	}
	return filepath.Abs(path.Clean(os.ExpandEnv(p)))
}

func WriteFile(file string, data []byte) error {
	expanded, err := ExpandPath(file)
	if err != nil {
		return err
	}
	if Exists(expanded) {
		info, err := os.Stat(expanded)
		if err != nil {
			return err
		}
		if info.Mode() != params.CurieIoConfig().ReadWriteExecutePermissions {
			return errors.New("file already exists without proper 0600 permissions")
		}
	}
	return os.WriteFile(expanded, data, params.CurieIoConfig().ReadWritePermissions)
}

// 환경 변수를 통해 Home 디렉토리 판단하는 함수
func HomeDir() string {
	// 'HOME'이라는 환경 변수가 설정되어 있으면
	// 그 환경 변수 그대로 반환
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

// 지정된 경로에 파일이 존재하는 경우 true 반환, 아니면 false 반환
func Exists(filename string) bool {
	filePath, err := ExpandPath(filename)
	if err != nil {
		return false
	}
	info, err := os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.WithError(err).Info("Checking for file existence returned an error")
		}
		return false
	}
	return info != nil && !info.IsDir()
}

func FileToBytes(file multipart.File) ([]byte, error) {
	var buf bytes.Buffer

	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
