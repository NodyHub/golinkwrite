# golinkwrite
Create a tar archive containing a provided file and a symlink that points to the write destination.

## Install

### Build from source

```shell
make install
```

### Fetch latest from GitHub

```shell
GOPRIVATE=github.com/NodyHub/golinkwrite go install github.com/NodyHub/golinkwrite@latest
```



## Usage

### Help

```shell
$ golinkwrite
Usage: golinkwrite <input> <target> <output> [flags]

Create a tar archive containing a symbolic link to a provided target and a provided file.

Arguments:
  <input>     Input file.
  <target>    Target destination in the filesystem.
  <output>    Output file.

Flags:
  -h, --help       Show context-sensitive help.
  -v, --verbose    Enable verbose output.
```

### In Action

```shell
$ echo 'Hello Alice :wave:!' | tee rabbit_hole.txt
Hello Alice :wave:!

$ golinkwrite -v rabbit_hole.txt /tmp/hi.txt alice.tar
time=2024-08-09T19:11:35.266+02:00 level=DEBUG msg="command line  parameters" cli="{Input:rabbit_hole.txt Target:/tmp/hi.txt Output:alice.tar Verbose:true}"
time=2024-08-09T19:11:35.266+02:00 level=DEBUG msg="input permissions" perm=-rw-r--r--
time=2024-08-09T19:11:35.266+02:00 level=DEBUG msg="input size" size=20
time=2024-08-09T19:11:35.266+02:00 level=INFO msg="tar file created" output=alice.tar

$ tar ztvf alice.tar
lrw-r--r--  0 0      0           0  1 Jan  1970 rabbit_hole.txt -> /tmp/hi.txt
-rw-r--r--  0 0      0          20  1 Jan  1970 rabbit_hole.txt
```
