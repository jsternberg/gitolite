// vim: set noexpandtab:
package gitolite

// ragel -Z -G2 shell.rl

import (
	"fmt"
)

%%{
	machine shell;

	action beg_str {
		beg = p
	}

	action end_str {
		buf = append(buf, data[beg:p]...)
	}

	action add_arg {
		args = append(args, string(buf))
		buf = make([]byte, 0)
	}

	chars = alnum | (punct - ["']);
	quote = "'" ((chars | '"') >beg_str %end_str)* "'";
	expr = ((chars+ >beg_str %end_str) | quote)+ %add_arg;
	main := space* expr (space+ expr)* space*;

	write data;
}%%

func shellArgs(data []byte) ([]string, error) {
	var cs, beg int
	var p, pe, eof int = 0, len(data), len(data)
	args := make([]string, 0)
	buf := make([]byte, 0)
%%{
	write init;
	write exec;
}%%

	if cs == shell_error {
		return nil, fmt.Errorf("an unknown parser error occurred")
	}

	if p != pe {
		return nil, fmt.Errorf("did not parse all input")
	}
	return args, nil
}
