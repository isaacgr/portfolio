/*
 * File: main.go
 * Author: Isaac Gluchowski-Rowell
 * Date Created: December 26, 2025
 *
 * Description: This file contains the main entry point of the application
 *              and demonstrates basic program flow.
 */

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isaacgr/loggir"
	"github.com/isaacgr/portfolio/internal/blog"
	"github.com/isaacgr/portfolio/internal/config"
	"github.com/isaacgr/portfolio/internal/webserver"
)

type Flags struct {
	ConfigFile string
	ConfigDir  string
	debug      bool
}

func parseFlags() Flags {

	var configDir = flag.String(
		"config-dir",
		"/var/lib/portfolio/",
		"Directory for config files",
	)

	var configFile = flag.String(
		"config",
		"config.ini",
		"Conifg file containing site initialization info",
	)

	var debug = flag.Bool(
		"debug",
		false,
		"Enable or disable debug logging",
	)

	flag.Parse()

	return Flags{
		ConfigFile: *configFile,
		ConfigDir:  *configDir,
		debug:      *debug,
	}
}

func main() {
	flags := parseFlags()

	debugLogging := flags.debug

	var log = loggir.GetLogger("portfolio", debugLogging)
	configLog := log.With("module", "config")
	portfolioLog := log.With("module", "api")
	webLog := log.With("module", "web")
	serverLog := log.With("module", "server")

	c, err := config.NewConfigProvider(
		flags.ConfigDir+flags.ConfigFile,
		configLog,
	)
	if err != nil {
		log.Error(
			"Unable to parse configuration. Exiting", "Error",
			err.Error(),
		)
		os.Exit(1)
	}

	host := c.Config.Host
	port := c.Config.Port
	blogPath := c.Config.BlogPath

	// TODO: I dont need to require that the blog exists to run, but that was
	// like 90% of this whole project
	if blogPath == "" {
		log.Error("Path to blogs not set. Exiting")
		os.Exit(1)
	}

	if _, err := os.ReadDir(blogPath); err != nil {
		log.Error("Unable to read 'blog_path' directory. Exiting")
		os.Exit(1)
	}

	recaptchaParams := webserver.ReCaptchaParams{
		Key:         c.Config.Key,
		ProjectID:   c.Config.ProjectID,
		Action:      c.Config.Action,
		Credentials: c.Config.Credentials,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	shutdownCtx, shutdownRelease := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	postChan := make(chan blog.Article, 1024)
	blogFinder := blog.NewBlogFinder(blogPath, []string{".md"}, nil, postChan, log)

	s := webserver.NewWebServer(
		host,
		port,
		shutdownCtx,
		webLog,
	)

	renderer := webserver.NewRenderer(
		log.With("webserver", "renderer"),
	)
	err = renderer.ParseTemplates()

	if err != nil {
		log.Error("Unable to parse templates.", "Error", err.Error())
	}

	responder := webserver.NewResponder()

	mailClient, err := webserver.NewMailClient(
		c.Config.SMTPHost,
		c.Config.SMTPUser,
		c.Config.SMTPPass,
		serverLog,
	)

	if err != nil {
		log.Error("Unable to create mail client.", "Error", err.Error())
	}

	portfolioApi := webserver.NewHandler(
		responder,
		portfolioLog,
		renderer,
		blogFinder,
		mailClient,
		recaptchaParams,
	)

	portfolioApi.Register(s.Handler)

	blogFinder.Start()
	s.Start()
	defer shutdownRelease()

	select {
	case <-sigChan:
		log.Info("Shutting down due to signal")
		s.Stop()
		signal.Stop(sigChan)
		close(sigChan)
	}
}
