package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"netinfo/internal/network"
	"netinfo/internal/preload"
	"netinfo/internal/receive"
	"netinfo/internal/send"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

func main() {
	// show mode
	var showPreload bool

	// send/receive mode
	var allowInsecure bool
	var encryptionKey string
	var interval time.Duration
	var endpoint string
	var username string
	var password string

	// send/receive mode file
	var filepath string

	// send/receive mode s3
	var regin string
	var stsToken string
	var pathStyle bool
	var bucket string
	var objectPath string

	// wireguard
	var remoteInterface string
	var wgInterface string
	var wgPeerKey string

	cmds := []*cli.Command{
		{
			Name:  "show",
			Usage: "show all network information",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "preload",
					Aliases:     []string{"p"},
					Usage:       "show preload information",
					Value:       false,
					Destination: &showPreload,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) (err error) {
				var bytes []byte
				var netInterfaces []network.NetInterface

				if showPreload {
					// 获取负载
					p, err := preload.NewPreload()
					if err != nil {
						return err
					}
					// 负载转换为比特流
					bytes, err = preload.Marshal(p, "json", nil)
					if err != nil {
						return err
					}
				} else {
					netInterfaces, err = network.GetNetInterfaces()
					bytes, err = json.Marshal(netInterfaces)
				}

				fmt.Println(string(bytes))
				return err
			},
		},
		{
			Name:    "send",
			Aliases: []string{"s"},
			Usage:   "send network information",

			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "send to filesystem",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "filepath",
							Usage:       "set file path",
							Required:    true,
							Destination: &filepath,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
						&cli.DurationFlag{
							Name:        "interval",
							Usage:       "set send interval",
							Destination: &interval,
						},
					},
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
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "allow_insecure",
							Usage:       "set allow insecure connect",
							Value:       false,
							Destination: &allowInsecure,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
						&cli.DurationFlag{
							Name:        "interval",
							Usage:       "set send interval",
							Destination: &interval,
						},
						&cli.StringFlag{
							Name:        "endpoint",
							Usage:       "set s3 server endpoint",
							Required:    true,
							Destination: &endpoint,
						},
						&cli.StringFlag{
							Name:        "regin",
							Usage:       "set s3 server regin",
							Value:       "us-east-1",
							Destination: &regin,
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
					},
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							send.ToS3Loop(endpoint, regin, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey), interval)
						} else {
							_, err = send.ToS3(endpoint, regin, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey))
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
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "allow_insecure",
							Usage:       "set allow insecure connect",
							Value:       false,
							Destination: &allowInsecure,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
						&cli.DurationFlag{
							Name:        "interval",
							Usage:       "set send interval",
							Destination: &interval,
						},
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
					},
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
			Usage:   "receive wireguard endpoint from network information",
			Flags: []cli.Flag{
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
				},
				&cli.DurationFlag{
					Name:        "interval",
					Aliases:     []string{"i"},
					Usage:       "set send interval",
					Destination: &interval,
				},
			},

			Commands: []*cli.Command{
				{
					Name:  "file",
					Usage: "receive network information from filesystem",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "filepath",
							Aliases:     []string{"f"},
							Usage:       "set file path",
							Required:    true,
							Destination: &filepath,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Aliases:     []string{"e"},
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							receive.FromFileLoop(filepath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey, interval)
						} else {
							err := receive.FromFile(filepath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey)
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
				{
					Name:  "s3",
					Usage: "receive network information from s3 server",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "allow_insecure",
							Usage:       "set allow insecure connect",
							Value:       false,
							Destination: &allowInsecure,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Aliases:     []string{"e"},
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
						&cli.StringFlag{
							Name:        "endpoint",
							Usage:       "set s3 server endpoint",
							Required:    true,
							Destination: &endpoint,
						},
						&cli.StringFlag{
							Name:        "regin",
							Usage:       "set s3 server regin",
							Value:       "us-east-1",
							Destination: &regin,
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
					},
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							receive.FromS3Loop(endpoint, regin, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey, interval)
						} else {
							err = receive.FromS3(endpoint, regin, username, password, stsToken, pathStyle, allowInsecure, bucket, objectPath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey)
							if err != nil {
								return err
							}
						}
						return nil
					},
				},
				{
					Name:  "webdav",
					Usage: "receive network information from webdav server",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "allow_insecure",
							Usage:       "set allow insecure connect",
							Value:       false,
							Destination: &allowInsecure,
						},
						&cli.StringFlag{
							Name:        "encryption_key",
							Aliases:     []string{"e"},
							Usage:       "set file encryption key",
							Destination: &encryptionKey,
						},
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
							Aliases:     []string{"f"},
							Usage:       "set webdav server filepath",
							Required:    true,
							Destination: &filepath,
						},
					},
					Action: func(ctx context.Context, cmd *cli.Command) (err error) {
						if interval != 0 {
							receive.FromWebDAVLoop(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey, interval)
						} else {
							err = receive.FromWebDAV(endpoint, username, password, allowInsecure, filepath, []byte(encryptionKey), remoteInterface, wgInterface, wgPeerKey)
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
			Name:    "decrypt",
			Aliases: []string{"d"},
			Usage:   "decrypt a file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "filepath",
					Aliases:     []string{"f"},
					Usage:       "set file path",
					Required:    true,
					Destination: &filepath,
				},
				&cli.StringFlag{
					Name:        "encryption_key",
					Aliases:     []string{"e"},
					Usage:       "set file encryption key",
					Required:    true,
					Destination: &encryptionKey,
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) (err error) {
				bytes, err := os.ReadFile(filepath)
				if err != nil {
					return err
				}
				plaintext, err := preload.Decrypt(bytes, []byte(encryptionKey))
				if err != nil {
					return err
				}
				fmt.Println(string(plaintext))
				return nil
			},
		},
	}

	// 打印版本函数
	cli.VersionPrinter = func(cmd *cli.Command) {
		fmt.Printf("%s\n", cmd.Root().Version)
	}

	cmd := &cli.Command{
		Usage:    "Network information manager",
		Version:  "v3.30",
		Commands: cmds,
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
