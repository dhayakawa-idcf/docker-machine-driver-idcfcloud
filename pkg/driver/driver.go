package driver

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/engine"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/state"
)

type Driver struct {
	*drivers.BaseDriver
	EnginePort int
	SSHKey     string
}

const (
	defaultTimeout = 15 * time.Second
)

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			Name:   "instance-type",
			Usage:  "Instance Type",
			Value:  "dummy",
			EnvVar: "INSTANCE_TYPE",
		},
	}
}

// NewDriver creates and returns a new instance of the driver
func NewDriver() drivers.Driver {
	return &Driver{
		EnginePort: engine.DefaultPort,
		BaseDriver: &drivers.BaseDriver{
			MachineName: "dummy",
			StorePath:   "dummy",
		},
	}
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "generic"
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHUsername() string {
	return d.SSHUser
}

func (d *Driver) GetSSHKeyPath() string {
	return d.SSHKeyPath
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.EnginePort = flags.Int("generic-engine-port")
	d.IPAddress = flags.String("generic-ip-address")
	d.SSHUser = flags.String("generic-ssh-user")
	d.SSHKey = flags.String("generic-ssh-key")
	d.SSHPort = flags.Int("generic-ssh-port")

	if d.IPAddress == "" {
		return errors.New("generic driver requires the --generic-ip-address option")
	}

	return nil
}

func (d *Driver) PreCreateCheck() error {
	if d.SSHKey != "" {
		if _, err := os.Stat(d.SSHKey); os.IsNotExist(err) {
			return fmt.Errorf("SSH key does not exist: %q", d.SSHKey)
		}

		// TODO: validate the key is a valid key
	}

	return nil
}

func (d *Driver) Create() error {
	if d.SSHKey == "" {
		log.Info("No SSH key specified. Assuming an existing key at the default location.")
	} else {
		log.Info("Importing SSH key...")

		d.SSHKeyPath = d.ResolveStorePath(path.Base(d.SSHKey))
		if err := copySSHKey(d.SSHKey, d.SSHKeyPath); err != nil {
			return err
		}

		if err := copySSHKey(d.SSHKey+".pub", d.SSHKeyPath+".pub"); err != nil {
			log.Infof("Couldn't copy SSH public key : %s", err)
		}
	}

	log.Debugf("IP: %s", d.IPAddress)

	return nil
}

func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("tcp://%s", net.JoinHostPort(ip, strconv.Itoa(d.EnginePort))), nil
}

func (d *Driver) GetState() (state.State, error) {
	address := net.JoinHostPort(d.IPAddress, strconv.Itoa(d.SSHPort))

	_, err := net.DialTimeout("tcp", address, defaultTimeout)
	if err != nil {
		return state.Stopped, nil
	}

	return state.Running, nil
}

func (d *Driver) Start() error {
	return nil
}

func (d *Driver) Stop() error {
	return nil
}

func (d *Driver) Restart() error {
	_, err := drivers.RunSSHCommandFromDriver(d, "sudo shutdown -r now")
	return err
}

func (d *Driver) Kill() error {
	return nil
}

func (d *Driver) Remove() error {
	return nil
}

func copySSHKey(src, dst string) error {
	if err := mcnutils.CopyFile(src, dst); err != nil {
		return fmt.Errorf("unable to copy ssh key: %s", err)
	}

	if err := os.Chmod(dst, 0600); err != nil {
		return fmt.Errorf("unable to set permissions on the ssh key: %s", err)
	}

	return nil
}
