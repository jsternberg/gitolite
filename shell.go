
//line shell.rl:1
// vim: set noexpandtab:
package gitolite

// ragel -Z -G2 shell.rl

import (
	"fmt"
)


//line shell.go:14
const shell_start int = 1
const shell_first_final int = 4
const shell_error int = 0

const shell_en_main int = 1


//line shell.rl:32


func shellArgs(data []byte) ([]string, error) {
	var cs, beg int
	var p, pe, eof int = 0, len(data), len(data)
	args := make([]string, 0)
	buf := make([]byte, 0)

//line shell.go:31
	{
	cs = shell_start
	}

//line shell.go:36
	{
	if p == pe {
		goto _test_eof
	}
	switch cs {
	case 1:
		goto st_case_1
	case 0:
		goto st_case_0
	case 4:
		goto st_case_4
	case 5:
		goto st_case_5
	case 2:
		goto st_case_2
	case 3:
		goto st_case_3
	case 6:
		goto st_case_6
	}
	goto st_out
	st1:
		if p++; p == pe {
			goto _test_eof1
		}
	st_case_1:
		switch data[p] {
		case 32:
			goto st1
		case 33:
			goto tr2
		case 39:
			goto st2
		}
		switch {
		case data[p] > 13:
			if 35 <= data[p] && data[p] <= 126 {
				goto tr2
			}
		case data[p] >= 9:
			goto st1
		}
		goto st0
st_case_0:
	st0:
		cs = 0
		goto _out
tr2:
//line shell.rl:13

		beg = p
	
	goto st4
tr9:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
//line shell.rl:13

		beg = p
	
	goto st4
	st4:
		if p++; p == pe {
			goto _test_eof4
		}
	st_case_4:
//line shell.go:105
		switch data[p] {
		case 32:
			goto tr8
		case 33:
			goto tr9
		case 39:
			goto tr10
		}
		switch {
		case data[p] > 13:
			if 35 <= data[p] && data[p] <= 126 {
				goto tr9
			}
		case data[p] >= 9:
			goto tr8
		}
		goto st0
tr8:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
//line shell.rl:21

		args = append(args, string(buf))
		buf = make([]byte, 0)
	
	goto st5
tr12:
//line shell.rl:21

		args = append(args, string(buf))
		buf = make([]byte, 0)
	
	goto st5
	st5:
		if p++; p == pe {
			goto _test_eof5
		}
	st_case_5:
//line shell.go:146
		switch data[p] {
		case 32:
			goto st5
		case 33:
			goto tr2
		case 39:
			goto st2
		}
		switch {
		case data[p] > 13:
			if 35 <= data[p] && data[p] <= 126 {
				goto tr2
			}
		case data[p] >= 9:
			goto st5
		}
		goto st0
tr10:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
	goto st2
	st2:
		if p++; p == pe {
			goto _test_eof2
		}
	st_case_2:
//line shell.go:175
		if data[p] == 39 {
			goto st6
		}
		if 33 <= data[p] && data[p] <= 126 {
			goto tr4
		}
		goto st0
tr4:
//line shell.rl:13

		beg = p
	
	goto st3
tr6:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
//line shell.rl:13

		beg = p
	
	goto st3
	st3:
		if p++; p == pe {
			goto _test_eof3
		}
	st_case_3:
//line shell.go:204
		if data[p] == 39 {
			goto tr7
		}
		if 33 <= data[p] && data[p] <= 126 {
			goto tr6
		}
		goto st0
tr7:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
	goto st6
	st6:
		if p++; p == pe {
			goto _test_eof6
		}
	st_case_6:
//line shell.go:223
		switch data[p] {
		case 32:
			goto tr12
		case 33:
			goto tr2
		case 39:
			goto st2
		}
		switch {
		case data[p] > 13:
			if 35 <= data[p] && data[p] <= 126 {
				goto tr2
			}
		case data[p] >= 9:
			goto tr12
		}
		goto st0
	st_out:
	_test_eof1: cs = 1; goto _test_eof
	_test_eof4: cs = 4; goto _test_eof
	_test_eof5: cs = 5; goto _test_eof
	_test_eof2: cs = 2; goto _test_eof
	_test_eof3: cs = 3; goto _test_eof
	_test_eof6: cs = 6; goto _test_eof

	_test_eof: {}
	if p == eof {
		switch cs {
		case 6:
//line shell.rl:21

		args = append(args, string(buf))
		buf = make([]byte, 0)
	
		case 4:
//line shell.rl:17

		buf = append(buf, data[beg:p]...)
	
//line shell.rl:21

		args = append(args, string(buf))
		buf = make([]byte, 0)
	
//line shell.go:268
		}
	}

	_out: {}
	}

//line shell.rl:42


	if cs == shell_error {
		return nil, fmt.Errorf("an unknown parser error occurred")
	}

	if p != pe {
		return nil, fmt.Errorf("did not parse all input")
	}
	return args, nil
}
