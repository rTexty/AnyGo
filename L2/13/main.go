package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseFields(spec string) ([]int, error) {
	if spec == "" {
		return nil, errors.New("empty fields spec")
	}

	parts := strings.Split(spec, ",")
	seen := make(map[int]struct{})
	fields := make([]int, 0, len(parts))

	addField := func(v int) error {
		if v <= 0 {
			return fmt.Errorf("field must be >= 1: %d", v)
		}
		if _, ok := seen[v]; ok {
			return nil
		}
		seen[v] = struct{}{}
		fields = append(fields, v)
		return nil
	}

	for _, raw := range parts {
		part := strings.TrimSpace(raw)
		if part == "" {
			return nil, errors.New("empty token in fields spec")
		}

		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("invalid range: %q", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid range start %q: %w", bounds[0], err)
			}
			end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid range end %q: %w", bounds[1], err)
			}
			if start <= 0 || end <= 0 {
				return nil, fmt.Errorf("range bounds must be >= 1: %q", part)
			}
			if start > end {
				return nil, fmt.Errorf("range start > end: %q", part)
			}

			for i := start; i <= end; i++ {
				if err := addField(i); err != nil {
					return nil, err
				}
			}
			continue
		}

		v, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid field %q: %w", part, err)
		}
		if err := addField(v); err != nil {
			return nil, err
		}
	}

	sort.Ints(fields)
	return fields, nil
}

func selectFields(line, delimiter string, fields []int, separated bool) (string, bool) {
	if !strings.Contains(line, delimiter) {
		if separated {
			return "", false
		}
		return line, true
	}

	parts := strings.Split(line, delimiter)
	out := make([]string, 0, len(fields))
	for _, idx := range fields {
		pos := idx - 1
		if pos < 0 || pos >= len(parts) {
			continue
		}
		out = append(out, parts[pos])
	}
	return strings.Join(out, delimiter), true
}

func main() {
	fieldsSpec := flag.String("f", "", "fields to print, e.g. 1,3-5")
	delimiter := flag.String("d", "\t", "delimiter character")
	separatedOnly := flag.Bool("s", false, "print only lines containing delimiter")
	flag.Parse()

	if *fieldsSpec == "" {
		fmt.Fprintln(os.Stderr, "usage: cut -f fields [-d delimiter] [-s]")
		os.Exit(2)
	}
	if *delimiter == "" {
		fmt.Fprintln(os.Stderr, "delimiter cannot be empty")
		os.Exit(2)
	}

	fields, err := parseFields(*fieldsSpec)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		res, ok := selectFields(line, *delimiter, fields, *separatedOnly)
		if !ok {
			continue
		}
		fmt.Println(res)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
