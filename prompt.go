package prompt

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var DefaultSelect = Select{
	Filter: true,
	Paging: 8,
}

func ReadLine(ctx context.Context, msg string) (string, error) {
	p := Prompt{
		Message: msg,
	}
	return p.ReadLine(ctx)
}

func Confirm(ctx context.Context, text string) (bool, error) {
	q := Question{
		Text: text,
	}
	return q.Confirm(ctx)
}

func SelectMultiple(ctx context.Context, msg string, opts []string) ([]int, error) {
	s := DefaultSelect
	s.Message = msg
	s.Options = opts
	return s.Multiple(ctx)
}

func SelectSingle(ctx context.Context, msg string, opts []string) (int, error) {
	s := DefaultSelect
	s.Message = msg
	s.Options = opts
	return s.Single(ctx)
}

type Prompt struct {
	Message string
	stdin   *bufio.Reader
}

func (p *Prompt) ReadLine(ctx context.Context) (string, error) {
	fmt.Print(p.Message)
	line, err := p.readLine(ctx)
	if err != nil {
		return "", err
	}
	return line, nil
}

func (p *Prompt) readLine(ctx context.Context) (string, error) {
	type lineAndError struct {
		line string
		err  error
	}
	if p.stdin == nil {
		p.stdin = bufio.NewReader(os.Stdin)
	} else {
		p.stdin.Reset(os.Stdin)
	}
	ch := make(chan lineAndError, 1)
	go func() {
		var m lineAndError
		m.line, m.err = p.stdin.ReadString('\n')
		m.line = strings.TrimSpace(m.line)
		ch <- m
	}()
	select {
	case m := <-ch:
		return m.line, m.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

type Question struct {
	Text string
}

func (q *Question) Confirm(ctx context.Context) (result bool, err error) {
	p := Prompt{
		Message: q.Text + " [Y/n]: ",
	}
	line, err := p.ReadLine(ctx)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(line, "y"), nil
}

type Select struct {
	Message string
	Options []string
	Filter  bool
	Paging  int

	index map[string]int
}

func (s *Select) init() {
	if s.index != nil {
		return
	}
	// Need to build an index to match selected options to their index.
	// This is required because of filtering and paging.
	s.index = make(map[string]int)
	for i, opt := range s.Options {
		s.index[opt] = i
	}
}

func (s *Select) Multiple(ctx context.Context) (answer []int, err error) {
	var p multiParser
	if err = s.doSelect(ctx, &p); err != nil {
		return nil, err
	}
	s.init()
	answer = make([]int, len(p.resp))
	for i, opt := range p.resp {
		var has bool
		answer[i], has = s.index[opt]
		if !has {
			panic("inconsistent state")
		}
	}
	return answer, nil
}

func (s *Select) Single(ctx context.Context) (answer int, err error) {
	var p singleParser
	if err = s.doSelect(ctx, &p); err != nil {
		return 0, err
	}
	s.init()
	answer, has := s.index[p.resp]
	if !has {
		panic("inconsistent state")
	}
	return answer, nil
}

type parser interface {
	pick([]string, string) error
	help([]string, func(pattern, desc string))
}

type singleParser struct {
	resp string
}

func (p *singleParser) pick(opts []string, line string) (err error) {
	i, err := strconv.Atoi(line)
	if err != nil {
		return err
	}
	if i >= len(opts) {
		return fmt.Errorf("index out of range")
	}
	p.resp = opts[i]
	return nil
}

func (p *singleParser) help(opts []string, fn func(pattern, desc string)) {
	max := "9"
	if n := len(opts); n > 0 && n < 10 {
		max = strconv.Itoa(n - 1)
	}
	fn("[0-"+max+"]", "Select option with given index")
}

type multiParser struct {
	resp []string
}

func (p *multiParser) pick(opts []string, line string) (err error) {
	if line == "*" {
		p.resp = opts
		return nil
	}
	var (
		ss = strings.Split(line, ",")
		xs = make([]string, len(ss))
	)
	for i, s := range ss {
		x, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return err
		}
		if x >= len(opts) {
			return fmt.Errorf("index out of range")
		}
		xs[i] = opts[x]
	}
	p.resp = xs
	return nil
}

func (p *multiParser) help(opts []string, fn func(pattern, desc string)) {
	max := "9"
	if n := len(opts); n > 0 && n < 10 {
		max = strconv.Itoa(n - 1)
	}
	fn("[0-"+max+"]", "Select options with given index(es)")
	fn("*", "Select all listed options")
}

func (s *Select) doSelect(ctx context.Context, parser parser) error {
	var options = s.Options
	var offset int
top:
	fmt.Println(s.Message)
	for i, j := offset, 0; i < len(options) && (s.Paging <= 0 || j < s.Paging); i, j = i+1, j+1 {
		opt := options[i]
		fmt.Printf("%d) %s\n", i, strings.TrimSpace(opt))
	}
	p := Prompt{
		Message: "Command (? for help): ",
	}
	for {
		line, err := p.ReadLine(ctx)
		if err != nil {
			return err
		}
		if line == "?" {
			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintln(tw, "  ?\tPrint this help message")
			if n := s.Paging; n > 0 {
				if offset > 0 {
					fmt.Fprintf(tw, "  p\tShow %d previous option(s)\n", n)
				}
				if rem := len(options) - offset; rem > n {
					fmt.Fprintf(tw, "  n\tShow %d next option(s)\n", n)
				}
			}
			if s.Filter {
				fmt.Fprintln(tw, "  /\tFilter options and repeat")
			}
			if n := s.Paging; n > 0 || s.Filter {
				fmt.Fprintln(tw, "  !\tReset view")
			}
			parser.help(options, func(pattern, desc string) {
				fmt.Fprintf(tw, "  %s\t%s\n", pattern, desc)
			})
			tw.Flush()

			continue
		}
		if s.Filter && line == "!" {
			options = s.Options
			offset = 0
			goto top
		}
		if s.Filter && line != "" && line[0] == '/' {
			options = filter(options, strings.TrimSpace(line[1:]))
			goto top
		}
		if n := s.Paging; n > 0 && line == "p" {
			if offset == 0 {
				fmt.Println("No options behind")
				continue
			}
			offset -= n
			goto top
		}
		if n := s.Paging; n > 0 && line == "n" {
			if rem := len(options) - offset; rem <= n {
				fmt.Println("No options ahead")
				continue
			}
			offset += n
			goto top
		}
		if err = parser.pick(options, line); err == nil {
			return nil
		}
		fmt.Printf("Unexpected input: %v\n", err)
	}
}

func filter(opts []string, pat string) (ret []string) {
	for _, opt := range opts {
		if match(opt, pat) {
			ret = append(ret, opt)
		}
	}
	return ret
}
