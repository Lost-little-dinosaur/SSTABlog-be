package Redis

import (
	"SSTABlog-be/internal/logger"
	"SSTABlog-be/internal/service/github"
	"encoding/json"
)

type GithubLogin struct {
	Info  github.GithubUserInfo
	State string
}

func NewGithubLoginString(m *github.GithubUserInfo, state string) []byte {
	var res GithubLogin
	res.State = state
	res.Info = *m
	marshal, _ := json.Marshal(res)
	return marshal
}

func RevertGithubLoginStruct(raw []byte) (res GithubLogin) {
	err := json.Unmarshal(raw, &res)
	if err != nil {
		logger.Error.Println(err)
	}
	return
}
