package main

func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := imageFile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <- errors; err != nil {
			return err
		}
	}

	return nil
}

func imageFile(f string) (string, error) {
	return "", nil
}