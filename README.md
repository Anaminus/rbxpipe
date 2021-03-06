### This project is archived! Pull requests will be ignored. Dependencies on this project should be avoided. Please fork this project if you wish to continue development.
----

# rbxpipe

**rbxpipe** enables Lua scripts to be piped into Roblox Studio, the output of
which can then be piped elsewhere.

```batch
echo print(BrickColor.Green().Color) | rbxpipe > brickcolor.txt

type brickcolor.txt
0.156863, 0.498039, 0.278431
```

## Options

Flags can have any of the following syntax:
```
-flag
-flag=x
-flag x
```
One or two minus signs may be used; they are equivalent. The last form is not
permitted for boolean flags because the meaning of the command `cmd -x *` will
change if there is a file called `0`, `false`, etc. You must use the
-flag=false form to turn off a boolean flag.

Flag parsing stops just before the first non-flag argument (`-` is a non-flag
argument) or after the terminator `--`. Supplying an unknown flag displays a
usage message.

More information on how flags are parsed can be found here:
http://golang.org/pkg/flag/

Option     | Description
-----------|------------
`-i`       | Specifies a Lua file to be executed. If unspecified, then data will be read from the standard input instead.
`-o`       | Specifies an output file. If unspecified, then the output will be written to the standard output instead.
`-studio`  | Specifies a path the studio executable. If unspecified, then rbxpipe will attempt to find it, assuming it is installed.
`-place`   | Specifies a Roblox place file to open with the script. If unspecified, then a new, empty place will be opened instead.
`-timeout` | Will terminate the studio process after the given duration. If less than 0, then the timeout is disabled. The duration is specified by an amount followed by a unit prefix (e.g. `30s` for 30 seconds, `5m` for 5 minutes, `1h` for 1 hour). Multiple units can be specified. Defaults to 30 seconds.
`-filter`  | Filters the output by message type. Each character in the filter string includes output messages of a certain type: `o` for regular output, `i` for info, `w` for warnings, and `e` for errors. Defaults to `oiwe`, or all message types.
`-format`  | Writes the output in a specified format. Acceptable formats are: `json` and `xml`. These formats can be suffixed with `i` to make the output more readable. A blank or invalid format outputs the raw data.
`-play`    | If given, the studio will run in `Play Solo` mode, which starts the RunService and inserts a character.

## Installation

To install rbxpipe from source to `$GOPATH/bin`:

1. [Install Go](http://golang.org/doc/install)
2. To run Go commands, make sure [Git](http://git-scm.com/downloads) is installed.
3. Using a shell that runs git (e.g. Git Bash), run the following command:
```sh
go get github.com/anaminus/rbxpipe
```

## Releases

Builds can be found on the [Releases](https://github.com/Anaminus/rbxpipe/releases) page.

Builds are cross-compiled with [goxc](https://github.com/laher/goxc).
