# dumpthread

Command line tool to save a [nostr](https://github.com/fiatjaf/nostr) thread on local filesystem, including all the associated events (replies, reactions, zaps, boosts).

Every event will be saved in a separated json file using its ID as the file name.

## Build and install

to build/install `dumpthread` you need [Go](https://go.dev/), version 1.19 or later.

you can build the executable via the provided `Makefile`:

```bash
make build
```

or can compile it by hand:

```bash
go build -o dumpthread main.go
```

to install into the default Go binary directory:

```bash
make install
```

or manually:

```bash
go install github.com/pedemonte/dumpthread
```

## Usage

```bash
NAME:
   dumpthread save - save a thread on filesystem

USAGE:
   dumpthread save [command options] [arguments...]

OPTIONS:
   --outdir value, -o value                             The path where to save the json events (default: ".")
   --eventid value, -e value                            The eventID of the first event of the thread
   --relay value, -r value [ --relay value, -r value ]  A relays to use (can be use multiple times)
   --help, -h                                           show help
```

example:

```bash
mkdir zzz

./dumpthread save \
-o zzz \
-e note13hgxlpqqpcmugyxey3q9wwfkytpu354anvyvswjamxhy7txplvhq98veh2 \
-r wss://nos.lol \
-r wss://relay.nostr.bg
```

will produce the following output in the directory `zzz`:

```bash
0d4b0623a1cf1ea70645e3be9df2f65bbd383b9ae3b8681ac69972db8a8f125a.json
2263328df3af3b889cd721fc36fd0daf05aebc5b24153882eb7a8b3f7315786f.json
31a0145b4ee9adc290bfbe8b886c156f966a6aac690fc5674cf4660df6e06390.json
37c9ffbf9beef9a38d628294c8e50b85977c14b78c0f7d20ea4439569dd23258.json
3b516c56462b9e140648839a51194082b4d4f4a7b3c17fd27e9e42e205abc30f.json
3dc58dd9d2354ba7c5bc58cf4d467a3c8c89b71e9e2cfbcd971eead4455a6b50.json
42b378f9a683662ab0f2cedf4befbe7676ef46bc4f32aa70949a1c1f1438381a.json
56eca61b49390053a94fc4c8babed39025103bef44bc7017981cdff02873e93b.json
6db92b7f08dbeba068dc28b692d149ece0a69115a63f1e46c45ca4cf7c2ac9d5.json
7a88511b78b90a44f0507d9529fbb246ad9f019c59e495a986fade587de09451.json
80c69493f0ecb7c153a06799baddd0b2911b53c0ef4dd69ef5b0e83f7263a91f.json
817021a6c6f2eda17a611d9a49624f67a47919530fb75832d09b548563e7f85c.json
861541bd4a3fad57cfa2a76c566c7793fa8b1ee6ec4248125559556e2be5c1ab.json
8d6e021ee8a81236029062b8db799443fa4250f83a845195207f731b93da877a.json
8dd06f84000e37c410d9244057393622c3c8d2bd9b08c83a5dd9ae4f2cc1fb2e.json
933978c6865b13e07b3cd839a86930695b5ce89f8876135260fe275f19507896.json
96a02c03fbea95e1af63529f59f9aaa2a177ad37ee0cbde3b52255a15ffc74ff.json
975f07e3e3a3240ac6aabc768a22e019fdbb391423792cbd85a602726472acd4.json
9b07e76ccaf1330115ba8f3ec4b107bdfe3528ff017b35efae63ef4389e9d6fc.json
c4cfbd301436546eab68cc6e3216925778cd0e01752e0505cf077f3254ea9017.json
d43afec237ef00eb8b6d8ac1063b8f338805b703e3d87e4af719ab518bf62414.json
e776ac379bc52b7dc0865cf2d81865d42889fc3ad178b043c6a3fffad675f3bf.json
f57155b46913dcfc78b6903f60ec9c1f87b29d10ae6d78d1597d827d99b31f89.json
f6f804f1c7ac867c57e2a3944f15a16d15c62b1ab6b81ba1b9e70832bd1c4e5d.json
```

## FAQ

* The code is very ugly, will you fix it?: I know it's ugly, I just needed this tool on the fly and maybe someone else could find it useful! maybe I'll fix it one day...
* What happens if a file is already there?: it will be overwritten;
* There will ever be an option to publish a saved thread on relays?: maybe, if I'll need it;
