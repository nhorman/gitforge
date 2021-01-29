module git-forge

go 1.15

replace github.com/ktrysmt/go-bitbucket => github.com/nhorman/go-bitbucket v0.0.0-20210128221002-2b73967d210c

require (
	github.com/bigkevmcd/go-configparser v0.0.0-20210106142102-909504547ead
	//github.com/ktrysmt/go-bitbucket v0.9.7
	github.com/ktrysmt/go-bitbucket v0.0.0-20210128221002-2b73967d210c
	github.com/motemen/go-gitconfig v0.0.0-20160409144229-d53da5028b75
	gopkg.in/ini.v1 v1.62.0
	gopkg.in/src-d/go-git.v4 v4.13.1
)
