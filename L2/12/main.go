package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type options struct {
	after      int
	before     int
	context    int
	countOnly  bool
	ignoreCase bool
	invert     bool
	fixed      bool
	withNumber bool
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func buildMatcher(pattern string, opts options) (func(string) bool, error) {
	if opts.fixed {
		needle := pattern
		if opts.ignoreCase {
			needle = strings.ToLower(needle)
		}

		return func(line string) bool {
			subject := line
			if opts.ignoreCase {
				subject = strings.ToLower(subject)
			}
			return strings.Contains(subject, needle)
		}, nil
	}

	expr := pattern
	if opts.ignoreCase {
		expr = "(?i)" + expr
	}
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return func(line string) bool {
		return re.MatchString(line)
	}, nil
}

func matchLines(lines []string, matcher func(string) bool, invert bool) []bool {
	matches := make([]bool, len(lines))
	for i, line := range lines {
		ok := matcher(line)
		if invert {
			ok = !ok
		}
		matches[i] = ok
	}
	return matches
}

func countMatches(matches []bool) int {
	count := 0
	for _, m := range matches {
		if m {
			count++
		}
	}
	return count
}

func markOutputLines(matches []bool, before, after int) []bool {
	out := make([]bool, len(matches))
	for i, matched := range matches {
		if !matched {
			continue
		}

		start := i - before
		if start < 0 {
			start = 0
		}
		end := i + after
		if end >= len(matches) {
			end = len(matches) - 1
		}
		for j := start; j <= end; j++ {
			out[j] = true
		}
	}
	return out
}

func printLines(lines []string, out []bool, withNumber bool, w io.Writer) {
	for i, line := range lines {
		if !out[i] {
			continue
		}
		if withNumber {
			fmt.Fprintf(w, "%d:%s\n", i+1, line)
			continue
		}
		fmt.Fprintln(w, line)
	}
}

func main() {
	opts := options{}
	flag.IntVar(&opts.after, "A", 0, "print N lines of trailing context")
	flag.IntVar(&opts.before, "B", 0, "print N lines of leading context")
	flag.IntVar(&opts.context, "C", 0, "print N lines of output context")
	flag.BoolVar(&opts.countOnly, "c", false, "print only a count of matching lines")
	flag.BoolVar(&opts.ignoreCase, "i", false, "ignore case distinctions")
	flag.BoolVar(&opts.invert, "v", false, "invert match")
	flag.BoolVar(&opts.fixed, "F", false, "interpret pattern as fixed string")
	flag.BoolVar(&opts.withNumber, "n", false, "print line number with output lines")
	flag.Parse()

	if opts.context > 0 {
		opts.after = opts.context
		opts.before = opts.context
	}

	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: grep [-A N] [-B N] [-C N] [-c] [-i] [-v] [-F] [-n] pattern [file]")
		os.Exit(2)
	}

	pattern := args[0]

	var (
		reader io.Reader = os.Stdin
		file   *os.File
		err    error
	)
	if len(args) == 2 {
		file, err = os.Open(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	lines, err := readLines(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	matcher, err := buildMatcher(pattern, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	matches := matchLines(lines, matcher, opts.invert)

	if opts.countOnly {
		fmt.Println(countMatches(matches))
		return
	}

	out := markOutputLines(matches, opts.before, opts.after)
	printLines(lines, out, opts.withNumber, os.Stdout)
}
