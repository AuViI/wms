package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

const (
	screenshotWidth  = 1920
	screenshotHeight = 1080
	screenshotFolder = renderFolder
)

/// renderPictures saves current view for application flags
/// parameterised into a png file to be accessed and printed.
///
/// Requires that folder `hfscc` is inside the current workding
/// directory when executing.
func renderPictures() {
	// render is "" by default, turned off
	if *render == "" {
		return
	}

	// local vars and display switch
	cwd, _ := os.Getwd()
	locations := strings.Split(*render, ",")
	absFolder := path.Join(cwd, screenshotFolder)

	Ok("rendering pictures")
	//Continue(fmt.Sprintf("location array: %v", locations))

	toForecast := func(loc string) string {
		return fmt.Sprintf("http://localhost%s/forecast/%s", *port, loc)
	}

	toFilename := func(loc string) string {
		file := fmt.Sprintf("cap_%s.png", loc)
		return fmt.Sprintf(path.Join(absFolder, file))
	}

	windowSize := fmt.Sprintf(
		"--window-size=%d,%d",
		screenshotWidth,
		screenshotHeight,
	)

	if _, staterr := os.Stat(absFolder); os.IsNotExist(staterr) {
		os.Mkdir(absFolder, os.ModePerm)
	}

	// iteration over locations
	for _, l := range locations {
		lt := strings.TrimSpace(l)
		if lt != "" {
			Continue(fmt.Sprintf("render %s", l))
			// using electron application which is submodule
			// inside "./hfscc", so this only works when the
			// current process is executed inside the repos'
			// directory
			cmd := exec.Command(
				"google-chrome-stable",
				"--headless",
				"--screenshot",
				windowSize,
				toForecast(l))
			cmd.Run()
			cmd.Wait()

			time.Sleep(2 * time.Second)

			cmdCopy := exec.Command(
				"mv",
				"screenshot.png",
				toFilename(l))
			cmdCopy.Run()
			cmdCopy.Wait()

			// sleep so cache can recover and windows don't
			// overlap in pictures [not tested]
			time.Sleep(10 * time.Second)
		}
	}
	Ok("finish rendering pictures")
}
