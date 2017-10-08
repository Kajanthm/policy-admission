/*
Copyright 2017 Home Office All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:    "policy-admission",
		Author:  "Rohith Jayawardene",
		Email:   "gambol99@gmail.com",
		Usage:   "is a service used to enforce secuirty policies within a cluster",
		Version: fmt.Sprintf("%s (git+sha: %s)", Version, GitSHA),

		OnUsageError: func(context *cli.Context, err error, isSubcommand bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options, %s\n", err)
			return err
		},

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "listen",
				Usage:  "the network interace the service should listen on `INTERFACE`",
				Value:  ":8443",
				EnvVar: "LISTEN",
			},
			cli.StringFlag{
				Name:   "tls-cert",
				Usage:  "the path to a file containing the tls certificate `PATH`",
				EnvVar: "TLS_CERT",
			},
			cli.StringFlag{
				Name:   "tls-key",
				Usage:  "the path to a file containing the tls key `PATH`",
				EnvVar: "TLS_KEY",
			},
			cli.StringFlag{
				Name:   "policies",
				Usage:  "the path to a file containing the security policies `PATH`",
				EnvVar: "POLICIES",
			},
		},

		Action: func(cx *cli.Context) error {
			config := &Config{
				Listen:   cx.String("listen"),
				Policies: cx.String("policies"),
				TLSCert:  cx.String("tls-cert"),
				TLSKey:   cx.String("tls-key"),
			}

			// @step create the controller
			ctl, err := newAdmissionController(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] unable to initialize controller, %q", err)
				os.Exit(1)
			}

			// @step: start the service
			if err := ctl.startController(); err != nil {
				fmt.Fprintf(os.Stderr, "[error] unable to start controller, %q", err)
				os.Exit(1)
			}

			// @step setup the termination signals
			signalChannel := make(chan os.Signal)
			signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			<-signalChannel

			return nil
		},
	}

	app.Run(os.Args)
}