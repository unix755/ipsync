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
	var Key string
	// 是否允许不安全连接(tls 自签证书)
	var skipTLSVerify bool
	// 启用循环, 每一次运行之前的时间间隔
	var interval time.Duration

	// file 模式下本地存储文件地址, webdav 模式下服务端存储文件地址
	var path string

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
				Name:        "key",
				Aliases:     []string{"k"},
				Usage:       "key used for encryption",
				Value:       "",
				Sources:     cli.EnvVars("KEY"),
				Destination: &Key,
			},
			&cli.BoolFlag{
				Name:        "skipTLSVerify",
				Aliases:     []string{"s"},
				Usage:       "skip TLS verification in server connection",
				Value:       false,
				Sources:     cli.EnvVars("SKIP_TLS_VERIFY"),
				Destination: &skipTLSVerify,
			},
			&cli.DurationFlag{
				Name:        "interval",
				Aliases:     []string{"i"},
				Usage:       "interval between repetitive tasks",
				Value:       0 * time.Second,
				Sources:     cli.EnvVars("INTERVAL"),
				Destination: &interval,
			},
		}

		FileFlags = []cli.Flag{
			&cli.StringFlag{
				Name:        "file_path",
				Usage:       "file path used in file protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("FILE_PATH"),
				Destination: &path,
			},
		}

		S3Flags = []cli.Flag{
			&cli.StringFlag{
				Name:        "s3_endpoint",
				Usage:       "endpoint used in s3 protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("S3_ENDPOINT"),
				Destination: &endpoint,
			},
			&cli.StringFlag{
				Name:        "s3_region",
				Usage:       "region used in s3 protocol",
				Value:       "us-east-1",
				Sources:     cli.EnvVars("S3_REGION"),
				Destination: &region,
			},
			&cli.StringFlag{
				Name:        "s3_access_key_id",
				Usage:       "access key id used in s3 protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("S3_ACCESS_KEY_ID"),
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "s3_secret_access_key",
				Usage:       "secret access key used in s3 protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("S3_SECRET_ACCESS_KEY"),
				Destination: &password,
			},
			&cli.StringFlag{
				Name:        "s3_sts_token",
				Usage:       "sts token used in s3 protocol",
				Value:       "",
				Sources:     cli.EnvVars("S3_STS_TOKEN"),
				Destination: &stsToken,
			},
			&cli.BoolFlag{
				Name:        "s3_path_style",
				Usage:       "path style used in s3 protocol",
				Value:       false,
				Sources:     cli.EnvVars("S3_PATH_STYLE"),
				Destination: &pathStyle,
			},
			&cli.StringFlag{
				Name:        "s3_bucket",
				Usage:       "bucket used in s3 protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("S3_BUCKET"),
				Destination: &bucket,
			},
			&cli.StringFlag{
				Name:        "s3_object_path",
				Usage:       "object path used in s3 protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("S3_OBJECT_PATH"),
				Destination: &objectPath,
			},
		}

		WebDAVFlags = []cli.Flag{
			&cli.StringFlag{
				Name:        "webdav_endpoint",
				Usage:       "endpoint used in webdav protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("WEBDAV_ENDPOINT"),
				Destination: &endpoint,
			},
			&cli.StringFlag{
				Name:        "webdav_username",
				Usage:       "username used in webdav protocol",
				Value:       "",
				Sources:     cli.EnvVars("WEBDAV_USERNAME"),
				Destination: &username,
			},
			&cli.StringFlag{
				Name:        "webdav_password",
				Usage:       "password used in webdav protocol",
				Value:       "",
				Sources:     cli.EnvVars("WEBDAV_PASSWORD"),
				Destination: &password,
			},
			&cli.StringFlag{
				Name:        "webdav_path",
				Usage:       "path used in webdav protocol",
				Required:    true,
				Value:       "",
				Sources:     cli.EnvVars("WEBDAV_PATH"),
				Destination: &path,
			},
		}
	)

	cmds := []*cli.Command{
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "send network information to a remote server as a file",
			Flags:   CommonFlags,
			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "send to local filesystem",
					Flags: FileFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							send.ToFileLoop(path, []byte(Key), interval)
						} else {
							return send.ToFile(path, []byte(Key))
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
							send.ToS3Loop(endpoint, region, username, password, stsToken, pathStyle, skipTLSVerify, bucket, objectPath, []byte(Key), interval)
						} else {
							_, err = send.ToS3(endpoint, region, username, password, stsToken, pathStyle, skipTLSVerify, bucket, objectPath, []byte(Key))
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
							send.ToWebDAVLoop(endpoint, username, password, skipTLSVerify, path, []byte(Key), interval)
						} else {
							_, err = send.ToWebDAV(endpoint, username, password, skipTLSVerify, path, []byte(Key))
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
			Name:    "print",
			Aliases: []string{"p"},
			Usage:   "print network information from a file or local host.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "path",
					Aliases:     []string{"p"},
					Usage:       "file path used for printing network information",
					Value:       "",
					Sources:     cli.EnvVars("PRINT_PATH"),
					Destination: &path,
				},
				&cli.StringFlag{
					Name:        "key",
					Aliases:     []string{"k"},
					Usage:       "key used for encryption",
					Value:       "",
					Sources:     cli.EnvVars("PRINT_KEY"),
					Destination: &Key,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) (err error) {
				// 提供路径则从文件中解码 preload 并打印, 不提供则打印本机 preload
				if path != "" {
					bytes, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					plaintext, err := preload.Decrypt(bytes, []byte(Key))
					if err != nil {
						return err
					}
					fmt.Println(string(plaintext))
				} else {
					p, err := preload.NewPreload()
					if err != nil {
						return err
					}
					bytes, err := preload.Marshal(p, "json", []byte(Key))
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
			Usage:   "watchdog used to change the WireGuard endpoint.",
			Flags: append(CommonFlags,
				&cli.StringFlag{
					Name:        "wg_remote_interface",
					Usage:       "remote interface used WireGuard watchdog",
					Required:    true,
					Value:       "",
					Sources:     cli.EnvVars("WG_REMOTE_INTERFACE"),
					Destination: &remoteInterface,
				},
				&cli.StringFlag{
					Name:        "wg_interface",
					Usage:       "interface used WireGuard watchdog",
					Required:    true,
					Value:       "",
					Sources:     cli.EnvVars("WG_INTERFACE"),
					Destination: &wgInterface,
				},
				&cli.StringFlag{
					Name:        "wg_peer_key",
					Usage:       "peer key used WireGuard watchdog",
					Value:       "",
					Sources:     cli.EnvVars("WG_PEER_KEY"),
					Destination: &wgPeerKey,
				}),
			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "receive from filesystem",
					Flags: FileFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromFile(path, []byte(Key))
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
						p, err := receive.FromS3(endpoint, region, username, password, stsToken, pathStyle, skipTLSVerify, bucket, objectPath, []byte(Key))
						if err != nil {
							return err
						}
						return p.UpdateWireGuardEndPoint(remoteInterface, wgInterface, wgPeerKey, -1, interval)
					},
				},
				{
					Name:  "webdav",
					Usage: "receive from webdav server",
					Flags: WebDAVFlags,
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						// 获取 preload
						p, err := receive.FromWebDAV(endpoint, username, password, skipTLSVerify, path, []byte(Key))
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
		Usage:    "IP Sync Tool from https://github.com/unix755/ipsync",
		Version:  "v3.33",
		Commands: cmds,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
