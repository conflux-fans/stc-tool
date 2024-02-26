package db

import "fmt"

func KeyFileName(name string) string {
	return fmt.Sprintf("FILE-%s", name)
}
