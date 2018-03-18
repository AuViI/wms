package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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
	locations := strings.Split(*render, ",")
	os.Setenv("DISPLAY", ":0")

	// logging
	Ok("rendering pictures")
	Continue(fmt.Sprintf("location array: %v", locations))

	// iteration over locations
	for _, l := range locations {
		lt := strings.TrimSpace(l)
		if lt != "" {
			Continue(fmt.Sprintf("rendering picture for%s", l))
			// using electron application which is submodule
			// inside "./hfscc", so this only works when the
			// current process is executed inside the repos'
			// directory
			cmd := exec.Command("electron", "hfscc", l)
			cmd.Run()
			cmd.Wait()
			// sleep so cache can recover and windows don't
			// overlap in pictures [not tested]
			<-time.After(50 * time.Second)
		}
	}
	Ok("finish rendering pictures")
}
