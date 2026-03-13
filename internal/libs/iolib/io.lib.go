package iolib

import (
	"fmt"
	"homedy/internal/libs/slicelib"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CopyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents using io.Copy
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Ensure all data is written to the destination file
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func SetPermission(path string, permissions []int, recursive bool) error {
	if len(permissions) != 3 {
		return fmt.Errorf("%w: must be 3 items", ErrPermissionLength)
	}
	for _, permission := range permissions {
		if permission > 7 || permission < 0 {
			return fmt.Errorf("%w: %d", ErrPermissionNotKnown, permission)
		}
	}

	if !Exists(path) {
		return fmt.Errorf("%w: %s", os.ErrNotExist, path)
	}

	permissionSlcStr := slicelib.Map(permissions, func(i int, p int) string { return strconv.Itoa(p) })
	permissionStr := strings.Join(permissionSlcStr, "")
	permissionInt, err := strconv.ParseInt(permissionStr, 8, 32)
	if err != nil {
		return fmt.Errorf("failed to parse permission: %w", err)
	}

	fileMode := os.FileMode(permissionInt)

	if !recursive {
		return os.Chmod(path, fileMode)
	}

	return filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(p, fileMode)
	})
}

func MakeDirWithPerm(path string, permission []int) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("error while make dir all: %w", err)
	}
	err = SetPermission(path, permission, true)
	if err != nil {
		return fmt.Errorf("error while set permission: %w", err)
	}
	return nil
}
