package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mgutz/ansi"

	git "gopkg.in/src-d/go-git.v4"
	gitObject "gopkg.in/src-d/go-git.v4/plumbing/object"
)

var verbose = flag.Bool("v", false, "Set verbose logging")

type DepStatus struct {
	Name                 string
	CurrentHash          string
	DependencyHash       string
	LengthBetweenCommits int
	Found                bool
	CleanTree            bool
}

func main() {
	flag.Parse()

	gdeps, err := loadGodeps()
	if err != nil {
		fmt.Println("Error trying to read Godeps file: " + err.Error())
		os.Exit(1)
	}

	depStatuses := make([]DepStatus, 0)
	for _, dep := range gdeps.Deps {
		d, err := getDepStatus(dep)
		if err != nil {
			fmt.Println("WARNING: Error trying to get dependency status of " + dep.ImportPath + ": " + err.Error())
		} else {
			depStatuses = append(depStatuses, *d)
		}
	}

	printResults(depStatuses, *verbose)
}

func getDepStatus(dep Dependency) (*DepStatus, error) {
	depHash := dep.Rev
	depPath := dep.ImportPath

	path := os.Getenv("GOPATH")

	repo, err := git.PlainOpen(path + "/src/" + depPath)
	if err != nil {
		return nil, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := wt.Status()
	if err != nil {
		return nil, err
	}
	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	distance, found := getLengthBetweenCommits(repo, depHash)

	return &DepStatus{
		Name:                 depPath,
		DependencyHash:       depHash,
		CurrentHash:          headRef.Hash().String(),
		CleanTree:            status.IsClean(),
		LengthBetweenCommits: distance,
		Found:                found,
	}, nil
}

func getLengthBetweenCommits(repo *git.Repository, depHash string) (int, bool) {
	commits := make([]*gitObject.Commit, 0)
	cIter, _ := repo.CommitObjects()
	_ = cIter.ForEach(func(c *gitObject.Commit) error {
		commits = append(commits, c)
		return nil
	})
	gitObject.ReverseSortCommits(commits)

	count := 0
	found := false

	for _, c := range commits {
		if c.Hash.String() == depHash {
			found = true
			break
		}
		count++
	}

	if !found {
		return 0, found
	}
	return count, found
}

func printResults(depStatuses []DepStatus, verbose bool) {
	red := ansi.ColorFunc("red")
	yellow := ansi.ColorFunc("yellow")
	green := ansi.ColorFunc("green")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "Dependency status:")

	for _, d := range depStatuses {
		detailMessages := make([]interface{}, 0)
		msg := ""
		if !d.CleanTree {
			msg = red(d.Name)
			detailMessages = append(detailMessages, red("Tree is not clean!"))
		}
		if msg == "" && !d.Found {
			msg = red(d.Name)
			detailMessages = append(detailMessages, red("Commit not found. Tree is not updated."))
			detailMessages = append(detailMessages, " - Head \t "+d.CurrentHash)
			detailMessages = append(detailMessages, " - Dependency \t "+d.DependencyHash)
		}
		if msg == "" && d.Found && d.LengthBetweenCommits > 0 {
			msg = yellow(d.Name)
			detailMessages = append(detailMessages, yellow("OLD commit found. Newer version available."))
			detailMessages = append(detailMessages, fmt.Sprintf(" Head is forward of %d commits:", d.LengthBetweenCommits))
			detailMessages = append(detailMessages, " - Head \t "+d.CurrentHash)
			detailMessages = append(detailMessages, " - Dependency \t "+d.DependencyHash)
		}
		if msg == "" {
			msg = green(d.Name)
		}

		fmt.Fprintln(w, "-", msg)

		if verbose {
			for _, dm := range detailMessages {
				fmt.Fprintln(w, "\t", dm)
			}
		}
	}

	w.Flush()
}
