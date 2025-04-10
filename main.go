package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run command",
				Description: "Run command in similar way that MCP API works. Mostly it's for debug.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "cmd",
						Usage: "command to run",
					},
					&cli.StringSliceFlag{
						Name:  "args",
						Usage: "arguments list",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					bb, err := execute(ctx, c.String("cmd"), c.StringSlice("args")...)
					if err != nil {
						return err
					}

					fmt.Println(string(bb))

					return nil
				},
			},
			{
				Name:  "serve",
				Usage: "run mcp server",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:  "cmds",
						Usage: "commands list",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					s := server.NewMCPServer(
						"MPCsh 🤖",
						"0.1.0",
					)

					commands := make([]mcp.Tool, 0)

					for _, cmdName := range c.StringSlice("cmds") {
						split := strings.Split(cmdName, ":")
						if len(split) > 1 {
							tool := mcp.NewTool(split[0],
								mcp.WithDescription(split[1]),
								mcp.WithString("args",
									mcp.Description("args to execute command with"),
								),
							)

							commands = append(commands, tool)

							continue
						}

						bb, err := exec.Command("man", cmdName).CombinedOutput()
						if err != nil {
							log.Fatal(err)
						}

						tool := mcp.NewTool(cmdName,
							mcp.WithDescription(string(bb)),
							mcp.WithString("args",
								mcp.Description("args to execute command with"),
							),
						)

						commands = append(commands, tool)
					}

					for _, t := range commands {
						RegisterCMD(s, t)
					}

					RegisterAddTool(s)

					// Start the stdio server
					if err := server.ServeStdio(s); err != nil {
						fmt.Printf("Server error: %v\n", err)
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func RegisterCMD(s *server.MCPServer, tool mcp.Tool) {
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		argsP, ok := request.Params.Arguments["args"]
		if ok {
			args, ok := argsP.(string)
			if !ok {
				return nil, errors.New("args must be a string")
			}

			bb, err := execute(ctx, tool.Name, args)
			if err != nil {
				return nil, err
			}

			return mcp.NewToolResultText(string(bb)), nil
		}

		bb, err := execute(ctx, tool.Name)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(bb)), nil
	})
}

func RegisterAddTool(s *server.MCPServer) {
	tool := mcp.NewTool("add-tool",
		mcp.WithDescription("add console tool to MCP server"),
		mcp.WithString("cmd",
			mcp.Description("command to register"),
		),
		mcp.WithString("description",
			mcp.Description("description to use the tool"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cmd, ok := request.Params.Arguments["cmd"].(string)
		if !ok {
			return nil, errors.New("cmd must be a string")
		}

		desc, ok := request.Params.Arguments["description"].(string)
		if !ok {
			return nil, errors.New("description must be a string")
		}

		RegisterCMD(s, mcp.NewTool(cmd,
			mcp.WithDescription(desc),
			mcp.WithString("args",
				mcp.Description("args to execute command with"),
			),
		))

		return mcp.NewToolResultText("OK"), nil
	})
}

func execute(ctx context.Context, cmd string, args ...string) (string, error) {
	var res strings.Builder

	c := exec.CommandContext(ctx, cmd)

	if len(args) > 0 && args[0] != "" {
		list := make([]string, 0)

		for _, a := range args {
			args, err := shellquote.Split(a)
			if err != nil {
				return "", err
			}

			list = append(list, args...)
		}

		c = exec.CommandContext(ctx, cmd, list...)
	}

	defer c.Cancel()

	c.Stdout = &res
	c.Stderr = &res

	if err := c.Run(); err != nil {
		return res.String(), fmt.Errorf("%w: cannot run %s %v", err, cmd, args)
	}

	return res.String(), nil
}

