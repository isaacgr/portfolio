package git

type GitManager struct {
	Host     string
	Token    *string
	PollFreq *int
	Client   GitClient
}

func NewGitManager(host string, token string, pollFreq int) *GitManager {
	return &GitManager{
		Host:     host,
		Token:    &token,
		PollFreq: &pollFreq,
	}
}
