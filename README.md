# gitforge
utility for manipulating git forges (bitbucket/gitlab/github/etc) as a git subcommand


## Overview
gitforge is mean to provide a semi generic set of commands for working with git
forges through the git command line.  Its very early in development, but soon we
hope to have the ability to do the following a fairly human readable and generic
way:

* Adding/Removing forge defintions
* Forking Repositories
* Cloning Repositories
* Syncing branches between parent and child repositories
* Opening/Closing/Reviewing Pull requests

Obviously, you can do these things with any forge through the web ui, or via the
command line with some tool or another using that forges rest api, but gitforge
hopes to make those tools more generic by creating a driver model, wherein a
user can create a forge defition to a given url and associate it with a forge
type that this project will use to select the appropriate driver to translate
that into the correct REST api model.  Currently bitbucket is being used as our
proof of concept, but it should be expandable to github/gitlab/pagure/etc in
short order

## Building
go build

## Installing
copy git-forge to /usr/local/bin  the git porcelain command:
`git forge`
should then be available

Note there is also a man page that can be installed


