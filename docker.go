package building

func (b *B) Docker(args ...string) Tool {
	// $$$$ MAT check what happens with empty instructions
	return b.MakeTool("docker", "--version", "https://www.docker.com", "")
}
