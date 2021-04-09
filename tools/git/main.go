package main

import (
	"flag"
	"log"

	"github.com/argoproj/argo-cd/util/git"
)

var (
	gitUser     = flag.String("user", "", "git user")
	gitPassword = flag.String("password", "", "git user password")
	localDir    = flag.String("dir", "", "repo local dir")
	repo        = flag.String("repo", "", "git repo")
	cli         git.Client
	err         error
)

func init() {
	flag.Parse()

	factory := git.NewFactory()
	cli, err = factory.NewClient(*repo, *localDir, git.NewHTTPSCreds(*gitUser, *gitPassword, true), true)
	if err != nil {
		panic(err)
	}
}

func main() {
	initLocalRepoDir()
	fetchRepo()
}

func initLocalRepoDir() {
	err = cli.Init()
	if err != nil {
		log.Fatalf("init local repo  dir error:%v\n", err)
	}

}

func fetchRepo() {
	err = cli.Fetch()
	if err != nil {
		log.Fatalf("fetch repo error:%v\n", err)
	}
}
