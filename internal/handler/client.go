package handler

import (
    "goo/internal/client"
    "goo/internal/svc"
    "github.com/urfave/cli"
)

var Commands = []cli.Command{
    {
        Name:    "once",
        Aliases: []string{"q"},
        Usage:   "一次性脚本",
        Subcommands: []cli.Command{
            {
                Name:   "rename-key",
                Usage:  "重新设置redis key",
                Action: client.NewUser(svc.NewServiceContext()).RenameKey,
            },
        },
    },
}
