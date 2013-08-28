package mailing

import (
	"errors"
	"fmt"
	"github.com/kch42/simpleconf"
	"os/exec"
)

// SendmailMailer is a Mailer implementation that sends mail using a sendmail-ish executable.
//
// Sendmail-ish in this context means:
//
// 	* The executable must accept the from address using the `-f <addr>` switch
// 	* The executable must accept the to address as a parameter
// 	* The executable expects a mail (including all headers) on stdin that is terminated by EOF
type SendmailMailer struct {
	Exec string
	Args []string
}

func (sm SendmailMailer) Mail(to, from string, msg []byte) (outerr error) {
	off := len(sm.Args)
	args := make([]string, off+3)
	copy(args, sm.Args)
	args[off] = "-f"
	args[off+1] = from
	args[off+2] = to

	cmd := exec.Command(sm.Exec, args...)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	defer func() {
		pipe.Close()
		err := cmd.Wait()
		if outerr == nil {
			outerr = err
		}
	}()

	_, err = pipe.Write(msg)
	return err
}

// SendmailMailerCreator creates an SendmailMailer using configuration values in the [mail] section.
//
// 	exec - The name of the executable
// 	argX - (optional) Additional arguments for the executable. X is a ascending number starting with 1
func SendmailMailerCreator(conf simpleconf.Config) (Mailer, error) {
	rv := SendmailMailer{}
	var err error

	if rv.Exec, err = conf.GetString("mail", "exec"); err != nil {
		return rv, errors.New("Missing [mail] exec config")
	}

	for i := 1; ; i++ {
		arg, err := conf.GetString("mail", fmt.Sprintf("arg%d", i))
		if err != nil {
			break
		}

		rv.Args = append(rv.Args, arg)
	}

	return rv, nil
}
