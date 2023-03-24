---
title: Authentication

weight: 30
---


{{<lead>}}
`hof mod` supports authenticating to code hosts
for private modules and increased rate limits.
{{</lead>}}

<br>

When you run `hof mod`, it looks for authentication methods
to use when fetching source code.
We recommend setting up one of these metheods for your primary code hosts,
especially GitHub for better rate limits in CI,
even if you only use public repositories.
The lookup order is:

1. .netrc
2. ENV variables
3. sshkey


### .netrc

Used if the `.netrc` file exists in the home directory and has a matching machine name.
On Windows, the file is `_netrc`.

{{<codeInner title=".netrc">}}
machine github.com
login github-token
password <github-token-value>
{{</codeInner>}}


### ENV variables

These environment variables are recognized for three major code hosts.

- `GITHUB_TOKEN`
- `GITLAB_TOKEN`
- `BITBUCKET_TOKEN`
- `BITBUCKET_USERNAME` and `BITBUCKET_PASSWORD`

### sshkey

The following locations are searched for sshkey configuration

- `~/.ssh/config`
- `/etc/ssh/config`
- `~/.ssh/id_rsa`

{{<codeInner title=".ssh/config">}}
Host github.com
	User git
	Hostname github.com
	PreferredAuthentications publickey
	IdentityFile /home/user/.ssh/id_rsa
{{</codeInner>}}
