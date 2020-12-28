package main

func main() {
	_, err := NewGitClient(".")
	if err != nil {
		return
	}
}
