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

	toForecast := func(loc string) string {
		return fmt.Sprintf("http://localhost%s/forecast/%s", *port, loc)
	}

	toDTage := func(loc string) string {
		return fmt.Sprintf("http://localhost%s/dtage/%s", *port, loc)
	}

	toFilename := func(loc, prefix string) string {
		file := fmt.Sprintf("cap_%s_%s.png", prefix, loc)
		return fmt.Sprintf(path.Join(absFolder, file))
	}

	windowSize := fmt.Sprintf(
		"--window-size=%d,%d",
		screenshotWidth,
		screenshotHeight,
	)

	fullRender := func(link, filename string) {
		Continue(fmt.Sprintf("render %s to %s", link, filename))
		cmd := exec.Command(
			"google-chrome-stable",
			"--headless",
			"--screenshot",
			windowSize,
			link)
		cmd.Run()
		cmd.Wait()

		time.Sleep(2 * time.Second)

		if _, staterr := os.Stat(absFolder); os.IsNotExist(staterr) {
			os.Mkdir(absFolder, os.ModePerm)
		}

		cmdCopy := exec.Command(
			"mv",
			"screenshot.png",
			filename)
		cmdCopy.Run()
		cmdCopy.Wait()

		// sleep so cache can recover and windows don't
		// overlap in pictures [not tested]
		time.Sleep(5 * time.Second)
	}

	// iteration over locations
	for _, l := range locations {
		lt := strings.TrimSpace(l)
		if lt != "" {
			fullRender(toForecast(lt), toFilename("forecast", lt))
			fullRender(toDTage(lt), toFilename("dtage", lt))
		}
	}
	Ok("finish rendering pictures")
}
