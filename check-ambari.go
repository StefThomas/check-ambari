package main

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/urfave/cli.v1"
	"os"
)

var debug bool
var ambariUrl string
var ambariLogin string
var ambariPassword string

func main() {

	// Logger setting
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Check service state on HDP/HDF cluster from Ambari API"
	app.Version = "develop"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-url",
			Usage:       "The ambari base URL",
			EnvVar:      "AMBARI_URL",
			Destination: &ambariUrl,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-login",
			Usage:       "The Ambari API login",
			EnvVar:      "AMBARI_LOGIN",
			Destination: &ambariLogin,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "ambari-password",
			Usage:       "The Ambari API password",
			EnvVar:      "AMBARI_PASSWORD",
			Destination: &ambariPassword,
		}),
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Display debug output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "check-service",
			Usage: "Check the service state",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name you should to check state",
				},
				cli.StringFlag{
					Name:  "service-name",
					Usage: "The service name you should to check state",
				},
				cli.BoolFlag{
					Name:  "exclude-node-alerts",
					Usage: "Use it if you should to exclude node alerts because of you already check by another way",
				},
			},
			Action: checkService,
		},
		{
			Name:  "check-node",
			Usage: "Check the node state",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "cluster-name",
					Usage: "The cluster name you should to check state",
				},
				cli.StringFlag{
					Name:  "node-name",
					Usage: "The node name you should to check",
				},
			},
			Action: checkNode,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Check the global parameter
func manageGlobalParameters() error {
	if debug == true {
		log.SetLevel(log.DebugLevel)
	}

	if ambariUrl == "" {
		return errors.New("You must set --ambari-url parameter")
	}

	if ambariLogin == "" {
		return errors.New("You must set --ambari-login parameter")
	}
	if ambariPassword == "" {
		return errors.New("You must set --ambari-password parameter")
	}

	return nil
}
