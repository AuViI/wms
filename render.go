package main

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
	secwait := 50
	locations := strings.Split(*render, ",")
	os.Setenv("DISPLAY", ":0")

	// logging
	fmt.Println("rendering pictures")
	fmt.Printf("working with: '%s'\non display: ':0'\n", *render)
	fmt.Printf("seconds between: %d\nlocation array: %v\n", secwait, locations)

	// iteration over locations
	for _, l := range locations {
		lt := strings.TrimSpace(l)
		if lt != "" {
			fmt.Println("rendering picture for", l)
			// using electron application which is submodule
			// inside "./hfscc", so this only works when the
			// current process is executed inside the repos'
			// directory
			cmd := exec.Command("electron", "hfscc", l)
			cmd.Run()
			cmd.Wait()
			// sleep so cache can recover and windows don't
			// overlap in pictures [not tested]
			<-time.After(secwait * time.Second)
		}
	}
	fmt.Println("finish rendering pictures")
}
