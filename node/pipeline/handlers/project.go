package handlers

import "strings"

func getProject(path string) string {
	return strings.ReplaceAll(path, "_btool.yml", "")
}
