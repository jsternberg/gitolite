package gitolite

import (
	"errors"
	"io"
	"os/exec"
	"path/filepath"
)

type UploadPack interface {
	UploadPack(path string, r io.ReadCloser, w io.WriteCloser) error
}

type ReceivePack interface {
	ReceivePack(path string, r io.ReadCloser, w io.WriteCloser) error
}

type Dir string

func (dir Dir) UploadPack(path string, r io.ReadCloser, w io.WriteCloser) error {
	path, err := dir.expandPath(path)
	if err != nil {
		return err
	}

	// This method has to be very careful with how it invokes the underlying process
	// related to the read and write closers.
	// In most programs, you would use the os/exec package and create a Cmd struct.
	// When you set Stdin and Stdout to the read and write objects, it should do the
	// correct thing in the normal circumstances.
	// But, we commonly pass in the channel itself to the read and write objects here.
	// If we pass in the channel as the Reader to Stdin, it will block forever waiting
	// for SSH to close the pipe (which SSH never does automatically).
	// So we need to manually spawn the process ourselves and be very careful with
	// how we interact with the pipes.
	cmd := exec.Command("git-upload-pack", path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// pipe the channel to stdin on a different channel
	errch := make(chan error, 1)
	go func() {
		_, err := io.Copy(stdin, r)
		stdin.Close()
		errch <- err
	}()

	_, err = io.Copy(w, stdout)
	if err != nil {
		// unable to copy stdout to the writer, abort this process and close the Reader
		r.Close()
		<-errch
		return err
	}

	state, err := cmd.Process.Wait()
	if err != nil || !state.Success() {
		// close the Reader since this has failed and SSH does not close this on its own
		r.Close()
		<-errch
		if err == nil {
			err = errors.New("process exited with bad status")
		}
		return err
	}
	return <-errch
}

func (dir Dir) expandPath(path string) (string, error) {
	return filepath.Join(string(dir), path), nil
}
