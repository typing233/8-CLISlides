package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/clislides/internal/executor"
	"github.com/user/clislides/internal/model"
	"github.com/user/clislides/internal/parser"
	"github.com/user/clislides/internal/server"
	"github.com/user/clislides/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "slides [file.md]",
	Short: "Terminal-based presentation tool",
	Long:  "A CLI tool for presenting Markdown slides in the terminal with full navigation and code execution support.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runPresentation,
}

var (
	sshFlag  bool
	portFlag int
	hostFlag string
)

func init() {
	rootCmd.Flags().BoolVar(&sshFlag, "ssh", false, "Start SSH server for remote sharing")
	rootCmd.Flags().IntVar(&portFlag, "port", 2222, "SSH server port")
	rootCmd.Flags().StringVar(&hostFlag, "host", "0.0.0.0", "SSH server host")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runPresentation(cmd *cobra.Command, args []string) error {
	content, err := loadContent(args)
	if err != nil {
		return fmt.Errorf("failed to load content: %w", err)
	}

	content = executor.Preprocess(content)

	pres, err := parser.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse slides: %w", err)
	}

	if overrideFromMeta(pres) {
		return server.Serve(pres, hostFlag, portFlag)
	}

	if sshFlag {
		return server.Serve(pres, hostFlag, portFlag)
	}

	return tui.Run(pres)
}

func loadContent(args []string) (string, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	return "", fmt.Errorf("no input file specified and no data on stdin\nUsage: slides <file.md>")
}

func overrideFromMeta(pres *model.Presentation) bool {
	if pres.Meta.SSHPort > 0 {
		portFlag = pres.Meta.SSHPort
		if pres.Meta.SSHHost != "" {
			hostFlag = pres.Meta.SSHHost
		}
		return true
	}
	return false
}
