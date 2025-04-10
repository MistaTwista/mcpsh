# mcpsh

Tool to execute any (allowed) cmd in OS.
This tool is PoC, just to check the tools.

## usage

```bash
./mcpsh --help
```

or

```bash
./mcpsh serve --cmds=ls --cmds="sh:run command in shell"
```

for `ls` description will be taken from `man ls`

### Using Docker

```bash
docker build -t mcpsh .
```

## inspector

```bash
npx @modelcontextprotocol/inspector ./mcpsh
# OR
npx @modelcontextprotocol/inspector docker run --rm -i mcpsh --cmds ls,apt-get
```
