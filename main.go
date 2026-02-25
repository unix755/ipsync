package main

import (
	"context"
	"fmt"
	"ipsync/internal/preload"
	"ipsync/internal/receive"
	"ipsync/internal/send"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	// 文件加密 key
	var encryptionKey string
	// 是否允许不安全连接(tls 自签证书)
	var allowInsecure bool
	// 启用循环, 每一次运行之前的时间间隔
	var interval time.Duration
	// 接收后保存到的文件
	var toFile string

	// file 模式下本地存储文件地址, webdav 模式下服务端存储文件地址
	var filepath string

	// s3 模式下服务器端点, webdav 模式下服务器端点
	var endpoint string
	// s3 模式下 access_key_id, webdav 模式下账号
	var username string
	// s3 模式下 secret_access_key, webdav 模式下密码
	var password string

	// s3 模式下独有变量
	var region string
	var stsToken string
	var pathStyle bool
	var bucket string
	var objectPath string

	// wireguard watchdog
	var remoteInterface string
	var wgInterface string
	var wgPeerKey string

	// flags
	var (
		CommonFlags = []cli.Flag{
			&cli.StringFlag{
				Name:        "encryption_key",
				Aliases:     []string{"k"},
				Usage:       "set encryption key",
				Destination: &encryptionKey,
			},
			&cli.BoolFlag{
				Name:        "allow_insecure",
				Usage:       "set allow insecure connect",
				Value:       false,
				Destination: &allowInsecure,
			},
			&cli.DurationFlag{
				Name:        "interval",
				Aliases:     []string{"i"},
				Usage:       "set interval",
				Destination: &interval,
			},
		}

		FileFlags = []cli.Flag{
			&cli.StringFlag{
				Name:        "filepath",
				Aliases:     []string{"f"},
				Usage:       "set file path",
				Required:    true,
				Destination: &filepath,
			},
		}

		S3Flags = []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Usage:       "set s3 server endpoint",
				Required:    true,
				Destination: &endpoint,
			},
			&cli.StringFlag{
				Name:        "region",
				Usage:       "set s3 server region",
				Value:       "us-east-1",
				Destination: &region,
			},
			&cli.StringFlag{
				Name:        "access_key_id",
				Usage:       "set s3 server access key id",
				Required:    true,
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "secret_access_key",
				Usage:       "set s3 server secret access key",
				Required:    true,
				Destination: &password,
			},
			&cli.StringFlag{
				Name:        "sts_token",
				Usage:       "set s3 server sts token",
				Destination: &stsToken,
			},
			&cli.BoolFlag{
				Name:        "path_style",
				Usage:       "set s3 server path style, false: virtual host, true: path",
				Value:       false,
				Destination: &pathStyle,
			},
			&cli.StringFlag{
				Name:        "bucket",
				Usage:       "set s3 server bucket",
				Required:    true,
				Destination: &bucket,
			},
			&cli.StringFlag{
				Name:        "object_path",
				Usage:       "set s3 server object path",
				Required:    true,
				Destination: &objectPath,
			},
		}

		WebDAVFlags = []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Usage:       "set webdav server endpoint",
				Required:    true,
				Destination: &endpoint,
			},
			&cli.StringFlag{
				Name:        "username",
				Usage:       "set webdav server username",
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "password",
				Usage:       "set webdav server password",
				Destination: &password,
			},
			&cli.StringFlag{
				Name:        "filepath",
				Usage:       "set webdav server filepath",
				Required:    true,
				Destination: &filepath,
			},
		}
	)

	cmds := []*cli.Command{
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "send network information",
			Flags:   CommonFlags,
			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "send to filesystem",
					Flags: FileFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							send.ToFileLoop(filepath, []byte(encryptionKey), interval)
						} else {
							return send.ToFile(filepath, []byte(encryptionKey))
						}
						return nil
					},
				},
				{
					Name:  "s3",
					Usage: "send to s3 server",
					Flags: S3Flags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							send.ToS3Loop(endpoint, region, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey), interval)
						} else {
							_, err = send.ToS3(endpoint, region, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey))
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
				{
					Name:  "webdav",
					Usage: "send to webdav server",
					Flags: WebDAVFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							send.ToWebDAVLoop(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey), interval)
						} else {
							_, err = send.ToWebDAV(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey))
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "receive",
			Aliases: []string{"r"},
			Usage:   "receive network information",
			Flags: append(CommonFlags,
				&cli.StringFlag{
					Name:        "to_file",
					Aliases:     []string{"t"},
					Usage:       "set save to file",
					Destination: &toFile,
				}),
			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "receive from filesystem",
					Flags: FileFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromFile(filepath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						// 保存为明文, 或者打印明文
						return p.SaveToFileOrPrint(toFile, []byte{})
					},
				},
				{
					Name:  "s3",
					Usage: "receive from s3 server",
					Flags: S3Flags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromS3(endpoint, region, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						// 保存为明文, 或者打印明文
						return p.SaveToFileOrPrint(toFile, []byte{})
					},
				},
				{
					Name:  "webdav",
					Usage: "receive network information from webdav server",
					Flags: WebDAVFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromWebDAV(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						// 保存为明文, 或者打印明文
						return p.SaveToFileOrPrint(toFile, []byte{})
					},
				},
			},
		},
		{
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "print network information",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "filepath",
					Aliases:     []string{"f"},
					Usage:       "set file path",
					Destination: &filepath,
				},
				&cli.StringFlag{
					Name:        "encryption_key",
					Aliases:     []string{"k"},
					Usage:       "set file encryption key",
					Destination: &encryptionKey,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) (err error) {
				if filepath != "" {
					bytes, err := os.ReadFile(filepath)
					if err != nil {
						return err
					}
					plaintext, err := preload.Decrypt(bytes, []byte(encryptionKey))
					if err != nil {
						return err
					}
					fmt.Println(string(plaintext))
				} else {
					p, err := preload.NewPreload()
					if err != nil {
						return err
					}
					bytes, err := preload.Marshal(p, "json", []byte(encryptionKey))
					if err != nil {
						return err
					}
					fmt.Println(string(bytes))
				}
				return nil
			},
		},
		{
			Name:    "watchdog",
			Aliases: []string{"w"},
			Usage:   "wireguard watchdog",
			Flags: append(CommonFlags,
				&cli.StringFlag{
					Name:        "remote_interface",
					Aliases:     []string{"r"},
					Usage:       "set remote interface",
					Required:    true,
					Destination: &remoteInterface,
				},
				&cli.StringFlag{
					Name:        "wg_interface",
					Aliases:     []string{"wi"},
					Usage:       "set wireguard interface",
					Required:    true,
					Destination: &wgInterface,
				},
				&cli.StringFlag{
					Name:        "wg_peer_key",
					Aliases:     []string{"wk"},
					Usage:       "set wireguard peer key",
					Destination: &wgPeerKey,
				}),
			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "receive from filesystem",
					Flags: FileFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromFile(filepath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						return p.UpdateWireGuardEndPoint(remoteInterface, wgInterface, wgPeerKey, -1, interval)
					},
				},
				{
					Name:  "s3",
					Usage: "receive from s3 server",
					Flags: S3Flags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromS3(endpoint, region, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						return p.UpdateWireGuardEndPoint(remoteInterface, wgInterface, wgPeerKey, -1, interval)
					},
				},
				{
					Name:  "webdav",
					Usage: "receive network information from webdav server",
					Flags: WebDAVFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromWebDAV(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey))
						if err != nil {
							return err
						}
						return p.UpdateWireGuardEndPoint(remoteInterface, wgInterface, wgPeerKey, -1, interval)
					},
				},
			},
		},
	}

	// 打印版本函数
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s\n", cmd.Root().Version)
	}

	cmd := &cli.Command{
		Usage:    "IP Sync Tool",
		Version:  "v3.31",
		Commands: cmds,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
