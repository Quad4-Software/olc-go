package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Quad4-Software/olc-go/olc"
)

func main() {
	decodeFlag := flag.String("decode", "", "decode this Open Location Code")
	shortenFlag := flag.Bool("shorten", false, "shorten encoded code relative to -ref-lat/-ref-lon")
	recoverFlag := flag.String("recover", "", "recover full code from short code using -ref-lat/-ref-lon")
	latStr := flag.String("lat", "-1.286386", "latitude decimal degrees")
	lonStr := flag.String("lon", "36.817223", "longitude decimal degrees")
	refLatStr := flag.String("ref-lat", "-1.286386", "reference latitude for shorten/recover")
	refLonStr := flag.String("ref-lon", "36.817223", "reference longitude for shorten/recover")
	codeLen := flag.Int("len", 10, "code length for encode")
	flag.Parse()

	if flag.NArg() > 0 {
		die(errors.New("unexpected arguments"))
	}

	if *decodeFlag != "" {
		area, err := olc.Decode(*decodeFlag)
		if err != nil {
			die(err)
		}
		clat, clng := area.Center()
		fmt.Printf("%s lo=(%f,%f) hi=(%f,%f) center=(%f,%f) len=%d\n",
			*decodeFlag, area.LatLo, area.LngLo, area.LatHi, area.LngHi, clat, clng, area.Len)
	}

	if *recoverFlag != "" {
		refLat, err := strconv.ParseFloat(*refLatStr, 64)
		if err != nil {
			die(err)
		}
		refLon, err := strconv.ParseFloat(*refLonStr, 64)
		if err != nil {
			die(err)
		}
		full, err := olc.RecoverNearest(*recoverFlag, refLat, refLon)
		if err != nil {
			die(err)
		}
		fmt.Printf("recover %s -> %s\n", *recoverFlag, full)
		return
	}

	lat, err := strconv.ParseFloat(*latStr, 64)
	if err != nil {
		die(err)
	}
	lon, err := strconv.ParseFloat(*lonStr, 64)
	if err != nil {
		die(err)
	}
	code := olc.Encode(lat, lon, *codeLen)
	fmt.Printf("encode %f %f len=%d -> %s\n", lat, lon, *codeLen, code)

	if *shortenFlag {
		refLat, err := strconv.ParseFloat(*refLatStr, 64)
		if err != nil {
			die(err)
		}
		refLon, err := strconv.ParseFloat(*refLonStr, 64)
		if err != nil {
			die(err)
		}
		short, err := olc.Shorten(code, refLat, refLon)
		if err != nil {
			die(err)
		}
		fmt.Printf("shorten -> %s\n", short)
	}
}

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
