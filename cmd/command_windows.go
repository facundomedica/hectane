package cmd

import (
	"github.com/hectane/hectane/cfg"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/util"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"path/filepath"
)

const (
	displayName = "Hectane"
	description = "Lightweight SMTP client"
)

// Run the specified command on the service.
func serviceCommand(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(exec.ServiceName)
	if err != nil {
		return err
	}
	defer s.Close()
	switch name {
	case "remove":
		return s.Delete()
	case "start":
		return s.Start()
	case "stop":
		_, err := s.Control(svc.Stop)
		return err
	}
	return nil
}

// Connect to the service manager and install the service. A default
// configuration is generated with the storage directory set to a folder in the
// same location as the executable.
var installCommand = &command{
	name:        "install",
	description: "install the service (Windows only)",
	exec: func() error {
		m, err := mgr.Connect()
		if err != nil {
			return err
		}
		defer m.Disconnect()
		exePath, err := util.Executable()
		if err != nil {
			return err
		}
		var (
			dir, _      = filepath.Split(exePath)
			cfgPath     = filepath.Join(dir, "config.json")
			storagePath = filepath.Join(dir, "storage")
			config      = &cfg.Config{}
		)
		config.Queue.Directory = storagePath
		if err := config.Save(cfgPath); err != nil {
			return err
		}
		if err := util.SecurePath(cfgPath); err != nil {
			return err
		}
		s, err := m.CreateService(exec.ServiceName, exePath, mgr.Config{
			StartType:   mgr.StartAutomatic,
			DisplayName: displayName,
			Description: description,
		}, "-config", cfgPath)
		if err != nil {
			return err
		}
		defer s.Close()
		if err := eventlog.InstallAsEventCreate(exec.ServiceName, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
			s.Delete()
			return err
		}
		return nil
	},
}

// Start the service.
var startCommand = &command{
	name:        "start",
	description: "start the service (Windows only)",
	exec: func() error {
		return serviceCommand("start")
	},
}

// Stop the service.
var stopCommand = &command{
	name:        "stop",
	description: "stop the service (Windows only)",
	exec: func() error {
		return serviceCommand("stop")
	},
}

// Remove the service.
var removeCommand = &command{
	name:        "remove",
	description: "remove the service (Windows only)",
	exec: func() error {
		if err := serviceCommand("remove"); err != nil {
			return err
		}
		return eventlog.Remove(exec.ServiceName)
	},
}

// Initialize the commands available for the current platform.
func Init() {
	commands = []*command{
		installCommand,
		removeCommand,
		startCommand,
		stopCommand,
	}
}
