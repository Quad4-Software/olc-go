package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"git.quad4.io/libraries/olc-go/olc"
)

func run(w io.Writer, args []string) error {
	fs := flag.NewFlagSet("example", flag.ContinueOnError)
	fs.SetOutput(w)
	decodeFlag := fs.String("decode", "", "decode this Open Location Code")
	shortenFlag := fs.Bool("shorten", false, "shorten encoded code relative to -ref-lat/-ref-lon")
	recoverFlag := fs.String("recover", "", "recover full code from short code using -ref-lat/-ref-lon")
	latStr := fs.String("lat", "-1.286386", "latitude decimal degrees")
	lonStr := fs.String("lon", "36.817223", "longitude decimal degrees")
	refLatStr := fs.String("ref-lat", "-1.286386", "reference latitude for shorten/recover")
	refLonStr := fs.String("ref-lon", "36.817223", "reference longitude for shorten/recover")
	codeLen := fs.Int("len", 10, "code length for encode")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() > 0 {
		return errors.New("unexpected arguments")
	}
	if *decodeFlag != "" && (*recoverFlag != "" || *shortenFlag) {
		return errors.New("-decode cannot be combined with -recover or -shorten")
	}
	if *recoverFlag != "" && *shortenFlag {
		return errors.New("-recover cannot be combined with -shorten")
	}

	parseCoordinate := func(flag, value string) (float64, error) {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid %s value %q: %w", flag, value, err)
		}
		return v, nil
	}
	reference := func() (float64, float64, error) {
		refLat, err := parseCoordinate("-ref-lat", *refLatStr)
		if err != nil {
			return 0, 0, err
		}
		refLon, err := parseCoordinate("-ref-lon", *refLonStr)
		if err != nil {
			return 0, 0, err
		}
		return refLat, refLon, nil
	}

	if *decodeFlag != "" {
		area, err := olc.Decode(*decodeFlag)
		if err != nil {
			return err
		}
		clat, clng := area.Center()
		fmt.Fprintf(w, "%s lo=(%f,%f) hi=(%f,%f) center=(%f,%f) len=%d\n",
			*decodeFlag, area.LatLo, area.LngLo, area.LatHi, area.LngHi, clat, clng, area.Len)
		return nil
	}

	if *recoverFlag != "" {
		refLat, refLon, err := reference()
		if err != nil {
			return err
		}
		full, err := olc.RecoverNearest(*recoverFlag, refLat, refLon)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "recover %s -> %s\n", *recoverFlag, full)
		return nil
	}

	lat, err := parseCoordinate("-lat", *latStr)
	if err != nil {
		return err
	}
	lon, err := parseCoordinate("-lon", *lonStr)
	if err != nil {
		return err
	}
	code := olc.Encode(lat, lon, *codeLen)
	fmt.Fprintf(w, "encode %f %f len=%d -> %s\n", lat, lon, *codeLen, code)

	if *shortenFlag {
		refLat, refLon, err := reference()
		if err != nil {
			return err
		}
		short, err := olc.Shorten(code, refLat, refLon)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "shorten -> %s\n", short)
	}
	return nil
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
