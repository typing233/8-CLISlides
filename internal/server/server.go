package server

import (
	"context"
	"fmt"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/clislides/internal/model"
	"github.com/user/clislides/internal/tui"
)

func Serve(pres *model.Presentation, host string, port int) error {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				m := tui.NewModel(pres)
				return m, []tea.ProgramOption{tea.WithAltScreen()}
			}),
		),
	)
	if err != nil {
		return err
	}

	fmt.Printf("SSH server listening on %s:%d\n", host, port)
	fmt.Println("Connect with: ssh localhost -p", port)
	return s.ListenAndServe()
}

func Shutdown(ctx context.Context) {
	_ = ctx
}
