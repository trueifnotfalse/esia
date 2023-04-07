package esia

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strings"
)

const defaultTmpPath = "/tmp"

type CliSigner struct {
	config *CliSignerConfig
}

func NewCliSigner(config *CliSignerConfig) *CliSigner {
	if 0 == len(config.TmpPath) {
		config.TmpPath = defaultTmpPath
	}
	return &CliSigner{
		config: config,
	}
}

func (s *CliSigner) Sign(message string) (string, error) {
	tmpInFilePath := s.getTmpFilePath()
	defer s.removeFile(tmpInFilePath)
	err := s.saveMessage(tmpInFilePath, message)
	if nil != err {
		return "", err
	}
	tmpOutFilePath := s.getTmpFilePath()
	defer s.removeFile(tmpOutFilePath)
	cmdArgs := s.createOpenSSLArgs(tmpInFilePath, tmpOutFilePath)
	err = s.runCmd(cmdArgs)
	if nil != err {
		return "", err
	}
	signed, err := os.ReadFile(tmpOutFilePath)
	if nil != err {
		println(err.Error())
		return "", err
	}

	return base64.URLEncoding.EncodeToString(signed), nil
}

func (s *CliSigner) removeFile(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if nil != err {
		return err
	}

	return os.RemoveAll(path)
}

func (s *CliSigner) runCmd(cmdArgs []string) error {
	cmd := exec.Command("openssl", cmdArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	stdout, err := cmd.Output()
	if nil != err {
		err = fmt.Errorf("error on sing message:'%v', stdout:'%s', stderr:'%s'", err, string(stdout), stderr.String())
		return err
	}

	return nil
}

func (s *CliSigner) createOpenSSLArgs(inFilePath, outFilePath string) []string {
	//password := ""
	//if 0 != len(s.config.PrivateKeyPassword) {
	//	password = " -passin pass:" + s.config.PrivateKeyPassword
	//}
	command := fmt.Sprintf(`smime -engine gost -sign -binary -outform DER -noattr -signer %s -inkey %s -in %s -out %s`, s.config.CertPath, s.config.PrivateKeyPath, inFilePath, outFilePath)
	return strings.Split(command, " ")
}

func (s *CliSigner) saveMessage(path, message string) error {
	err := os.WriteFile(path, []byte(message), os.ModePerm)
	if nil != err {
		return err
	}
	return nil
}

func (s *CliSigner) getTmpFilePath() string {
	return fmt.Sprintf("%s/%s", s.config.TmpPath, uuid.New().String())
}
