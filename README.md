# mpcsh

Tool to execute any (allowed) cmd in OS.
This tool is PoC, just to check the tools.

## usage

```bash
CMDS=ls,cat ./mpcsh
```

or

```bash
./mpcsh --cmds ls,cat
```

## inspector

```bash
npx @modelcontextprotocol/inspector ./mpcsh
```

then open inspector in browser and connect to server:

1. Transport: stdio
2. Command: ./mpcsh
3. Arguments: --cmds=ls,cat
