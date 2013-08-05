package models

import (
	"fmt"
	"os/exec"
	"bytes"
)
// A Cleaner is a function that manipulates bytes.
// It is used to clean up output from Docker.
type Cleaner func([]byte) []byte

// A command can execute a command and will run cleaners on the output.
type Command struct {
	Cmd      *exec.Cmd
	cleaners []Cleaner
}

// Create a command with a default cleaner of trimming space.
func NewCommand(cmd string, args ...string) *Command {
	return &Command{
		exec.Command(cmd, args...),
		[]Cleaner{bytes.TrimSpace},
	}
}

// Get the output and error from running the Command
// Also runs cleaners on the output
func (c *Command) Output() ([]byte, error) {
	out, err := c.Cmd.Output()
	for _, cleaner := range c.cleaners {
		out = cleaner(out)
	}
	return out, err
}

// Add a new cleaner for the Command to run after Output is called
func (c *Command) RegisterCleaner(cleaner func([]byte) []byte) {
	c.cleaners = append(c.cleaners, cleaner)
}

// Cleanup function for docker. Might get better with -cidfile.
func CleanUp() error {
	// Get a list of running processes
	cmd = NewCommand("docker", "ps", "-a", "-q")
	out, err = cmd.Output()
	if err != nil {
		return err
	}
	lines := bytes.Split(out, []byte("\n"))
	instances := make([]string, len(lines) + 1)
	// Hack because sometimes go's type system is less flexible than it seems
	instances[0] = "rm"
	for i, line := range lines {
		instances[i + 1] = string(line)
	}
	// This removes all running processes
	cmd = NewCommand("docker", instances...)
	out, err = cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

// A cleaner to only get the first line of output since piping via exec.Command is hard.
func FirstLine(in []byte) []byte {
	return bytes.Split(in, []byte("\n"))[0]
}

func main() {
	cmd := NewCommand("docker", "run", "7662e12d0778", "/dockulator/calculators/calc.rb", "23 + 44")
	out, _ := cmd.Output()
	fmt.Println("Output:", string(out))

	cmd = NewCommand("docker", "run", "7662e12d0778", "/bin/cat", "/etc/issue")
	cmd.RegisterCleaner(FirstLine)
	out, _ = cmd.Output()
	fmt.Println("OS:", string(out))
}


