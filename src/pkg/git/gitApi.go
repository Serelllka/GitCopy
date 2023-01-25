package gitApi

import (
	"CatFile/src/pkg/builder"
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"sync"
)

const (
	gitCommandName = "git"
	gitCatFile     = "cat-file"
)

var gitBinPath string

func init() {
	var err error
	gitBinPath, err = exec.LookPath(gitCommandName)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	_ = iota
	TreeObject
	BlobObject
	CommitObject
	TagObject
)

const (
	TreeName   = "tree"
	BlobName   = "blob"
	CommitName = "commit"
	TagName    = "tag"
)

type Object interface {
	GetHash() (string, error)
}

type Repository struct {
	mtx  sync.RWMutex
	path string
}

func NewRepository(path string) *Repository {
	return &Repository{
		path: path,
	}
}

func (r *Repository) getPath() (string, error) {
	if r.path == "" {
		return os.Getwd()
	}
	return r.path, nil
}

func (r *Repository) getGitCommand(buffer *bytes.Buffer, args []string) exec.Cmd {
	cmdBuilder := &builder.CommandBuilder{}
	cmdBuilder.SetPath(gitBinPath).SetStdout(buffer).SetArgs(args)
	return cmdBuilder.GetCommandInstance()
}

func (r *Repository) GetObjectType(hash string) (int, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	output := bytes.Buffer{}

	cmd := r.getGitCommand(&output, []string{gitBinPath, gitCatFile, "-t", hash})
	cmd.Dir = r.path

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	switch output.String() {
	case TreeName:
		return TreeObject, nil
	case BlobName:
		return BlobObject, nil
	case CommitName:
		return CommitObject, nil
	case TagName:
		return TagObject, nil
	default:
		return 0, errors.New(output.String())
	}
}

type Tree struct {
	hash string
}

type Blob struct {
	hash string
}

type Commit struct {
	hash string
}

type Tag struct {
	hash string
}
