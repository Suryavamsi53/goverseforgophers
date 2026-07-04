# Installing Go

To get started with Go, you first need to install it on your machine. You can always download the latest version of Go directly from the [official download page](https://go.dev/dl/).

Select the tab for your computer's operating system below, then follow its installation instructions.

## Linux

1. Download the latest `.linux-amd64.tar.gz` file.
2. Remove any previous Go installation by deleting the `/usr/local/go` folder (if it exists), then extract the archive you just downloaded into `/usr/local`, creating a fresh Go tree in `/usr/local/go`.

```bash
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.x.linux-amd64.tar.gz
```
*(You may need to run this command as root or through `sudo`)*

3. Add `/usr/local/go/bin` to the `PATH` environment variable. You can do this by adding the following line to your `$HOME/.profile` or `/etc/profile` (for a system-wide installation):

```bash
export PATH=$PATH:/usr/local/go/bin
```

4. Apply the changes immediately by running `source $HOME/.profile`.

## Mac

1. Download the latest `.pkg` file for macOS.
2. Open the package file you downloaded and follow the prompts to install Go.
3. The package installs the Go distribution to `/usr/local/go`. The package should automatically put the `/usr/local/go/bin` directory in your `PATH` environment variable. You may need to restart any open Terminal sessions for the change to take effect.

## Windows

1. Download the latest `.msi` file.
2. Open the MSI file you downloaded and follow the prompts to install Go.
3. By default, the installer will install Go to `Program Files` or `Program Files (x86)`. You can change the location as needed. After installing, you will need to close and reopen any open command prompts so that changes to the environment made by the installer are reflected at the command prompt.

## Verifying the Installation

Regardless of your operating system, verify that you've installed Go by opening a command prompt or terminal and typing the following command:

```bash
$ go version
```

Confirm that the command prints the installed version of Go. For example:
`go version go1.22.0 linux/amd64`

## Managing Go Installations

If you need to install multiple versions of Go (for example, to test your code against an older release), you can install them side-by-side. 

```bash
$ go install golang.org/dl/go1.20.7@latest
$ go1.20.7 download
$ go1.20.7 version
go version go1.20.7 linux/amd64
```
