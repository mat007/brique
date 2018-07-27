package b

func Docker(args ...string) Tool {
	// $$$$ MAT check what happens with empty instructions
	return MakeTool("docker", "--version", "https://www.docker.com", "")
}
