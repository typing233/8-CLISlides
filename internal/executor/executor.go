package executor

import (
	"os/exec"
	"regexp"
	"strings"
)

var preprocessorRe = regexp.MustCompile("(?m)^~~~([^\n]*)\n([\\s\\S]*?)^~~~$")

func Preprocess(content string) string {
	return preprocessorRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := preprocessorRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		command := strings.TrimSpace(sub[2])
		if command == "" {
			return match
		}
		cmd := exec.Command("sh", "-c", command)
		out, err := cmd.Output()
		if err != nil {
			return "```\nError: " + err.Error() + "\n```"
		}
		return "```\n" + strings.TrimRight(string(out), "\n") + "\n```"
	})
}

func ExecuteCodeBlock(code, lang string) (string, error) {
	var cmd *exec.Cmd
	switch lang {
	case "bash", "sh", "shell":
		cmd = exec.Command("sh", "-c", code)
	case "python", "python3":
		cmd = exec.Command("python3", "-c", code)
	case "go":
		cmd = exec.Command("go", "run", "-")
		cmd.Stdin = strings.NewReader(code)
	default:
		cmd = exec.Command("sh", "-c", code)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}
