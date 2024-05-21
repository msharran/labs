package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Example how to resolve a revision into its commit counterpart
func main() {
	CheckArgs("<path>", "<revision>")

	path := os.Args[1]
	revision := os.Args[2]

	// We instantiate a new repository targeting the given path (the .git     folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	h, err := r.ResolveRevision(plumbing.Revision(revision))

	CheckIfError(err)

	hash := h.String()

	refs, _ := r.References()
	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			if ref.Hash().String() == hash {
				hashRef := plumbing.NewHashReference(ref.Name(), ref.Hash())
				if hashRef.Name().IsBranch() {
					fmt.Println("Branch found", ref, hashRef.Name().Short())
					return nil
				}
			}
		}

		return nil
	})

}

func CheckIfError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CheckArgs(args ...string) {
	if len(os.Args) < len(args)+1 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
}
