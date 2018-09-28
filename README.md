# bosh tab completion

## Installation

Head over to the releases page and download the binary best suited to your
operating system. Then, you'll need to eval the command that prints out the
source code relevant to your shell.

What I mean is...

### For `bash` users

Add to your `.bashrc` or `.bash_profile`

```bash
eval "$(/path/to/bosh-complete bash-source)"
```

### For `zsh` users

Add to your `zshrc` or `.zprofile`

```zsh
eval "$(/path/to/bosh-complete zsh-source)"
```

`/path/to/bosh-complete` should be replaced with the location where you have
installed the bosh-complete binary.

## What If It Isn't Completing Something

First, make sure that you're logged into bosh on the target you're trying to
have it complete for - this tool reads your `.bosh/config` to determine auth
information to use. If that's out of date, then `bosh-complete` can't auth any
better than the bosh cli can (which is to say it cannot).

Also, be aware that completing some info is reliant upon you having already
provided the flag for some other piece of information. For example, the
`--deployment` flag can not be completed if the `--environment` flag has not
been given yet, because `bosh-complete` has no way of knowing at that point
which bosh director to query for deployment names!

Beyond that? I don't know! Maybe the thing you want completed isn't implemented
(yet). Maybe there's a bug (gasp!). Drop an issue on the repository and maybe we can get
through this together.