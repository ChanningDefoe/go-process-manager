package manager

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Process interface {
	Start() error
	Stop() error
}

type Manager struct {
	processes []Process
	signals   chan os.Signal
	errors    chan error
}

func NewManager(processes []Process) *Manager {
	return &Manager{
		processes: processes,
		signals:   make(chan os.Signal, 1),
		errors:    make(chan error, 1),
	}
}

func (m *Manager) Start() {
	signal.Notify(m.signals, syscall.SIGINT, syscall.SIGTERM)

	for _, process := range m.processes {
		go func(process Process) {
			if err := process.Start(); err != nil {
				m.errors <- err
			}
		}(process)
	}

	for {
		select {
		case err := <-m.errors:
			m.Stop()
			log.Fatalf("Error: %v", err)
		case <-m.signals:
			m.Stop()
			os.Exit(0)
		}
	}
}

func (m *Manager) Stop() {
	for _, process := range m.processes {
		if err := process.Stop(); err != nil {
			log.Printf("Error stopping process: %v", err)
		}
	}
}
