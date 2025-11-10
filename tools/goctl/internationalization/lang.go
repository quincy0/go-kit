package internationalization

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/go-github/v47/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var (
	// OwnerString github owner
	OwnerString string
	// RepoString github repo
	RepoString string
	// PathString repo,s path
	PathString string
	// MovePathString move,s path
	MovePathString string
)

//接口文档用于初始化生成国际化多语言目录文件
type InternationalInterface interface {
	ReadFile()
	MoveFileToProjectDir(path string)
}

type githubInternational struct {
	Owner, Repo, Path string
	Ctx               context.Context
	FileByte          map[string][]byte
}

func NewGithubInternational(owner, repo, path string) *githubInternational {
	return &githubInternational{
		Path:  path,
		Ctx:   context.Background(),
		Owner: owner,
		Repo:  repo,
	}
}

func (g *githubInternational) ReadFile() {
	//该token目前未员工私人创建永久不过期token
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "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"})
	tc := oauth2.NewClient(g.Ctx, ts)
	client := github.NewClient(tc)
	fileContent, directoryContent, resp, err := client.Repositories.GetContents(g.Ctx, g.Owner, g.Repo, g.Path, &github.RepositoryContentGetOptions{})
	if err != nil || resp.StatusCode != 200 {
		log.Fatalf("国际化内容获取异常，status: %v, err: %v", resp, err)
	}
	g.getAllContent(client, fileContent, directoryContent)
}

func (g *githubInternational) MoveFileToProjectDir(path string) {
	//判断目录
	_, err := os.Stat(path)
	if err != nil {
		//生成目录
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			log.Fatalf("国际化项目目录生成失败: err: %v", err)
		}
	}
	lastStr := path[len(path)-1:]
	for fileName, fileContent := range g.FileByte {
		fName := ""
		if lastStr != "/" {
			fName = fmt.Sprintf("%s/%s", path, fileName)
		} else {
			fName = fmt.Sprintf("%s%s", path, fileName)
		}

		f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
		if err != nil {
			log.Fatalf("写入文件失败： err: %v", err)
		}
		f.WriteString(string(fileContent))
		f.Close()
	}
}

func (g *githubInternational) getAllContent(client *github.Client, fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent) {
	g.FileByte = make(map[string][]byte, 0)
	if fileContent != nil {
		decodeBytes, _ := base64.StdEncoding.DecodeString(*fileContent.Content)
		g.FileByte[*fileContent.Name] = decodeBytes
	} else {
		for _, dirContent := range directoryContent {
			f, d, _, _ := client.Repositories.GetContents(g.Ctx, g.Owner, g.Repo, *dirContent.Path, &github.RepositoryContentGetOptions{})
			if len(d) > 0 {
				g.getAllContent(client, f, d)
			} else {
				decodeBytes, _ := base64.StdEncoding.DecodeString(*f.Content)
				g.FileByte[*dirContent.Name] = decodeBytes
			}
		}
	}
}

func MultilingualAction(_ *cobra.Command, _ []string) error {
	international := NewGithubInternational(OwnerString, RepoString, PathString)
	international.ReadFile()
	international.MoveFileToProjectDir(MovePathString)

	return nil
}
