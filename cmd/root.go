/*
* Copyright © 2024 Minand Manomohanan <minand.nell.mohan@gmail.com>
 */
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/minand-mohan/caching-proxy/server"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	port   string
	origin string
)

// rootCmd represents the base command when called without any subcommands
// This application works with the root command and switches functionality
// based on the arguments passed.
var rootCmd = &cobra.Command{
	Use:   "caching-proxy",
	Short: "A caching proxy server.",
	Long: `Caching Proxy Server is a high-performance, scalable proxy server designed to cache HTTP responses. 
It helps reduce latency and load on backend servers by serving cached responses for repeated requests. 
This server supports various features such as gzip compression handling, custom headers, and more, 
making it an ideal solution for improving the efficiency and speed of web applications.
Arguments:
  --port   Port on which the server will listen (required)
  --origin Origin URL to be proxied (required)`,
	Run: func(cmd *cobra.Command, args []string) {
		if port == "" || origin == "" {
			fmt.Println("port and origin must not be empty!")
			os.Exit(1)
		}
		// WaitGroup to wait for the goroutines to finish
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			server.SetUpServer(port, origin)
		}()

		wg.Wait()

		fmt.Println("Application stopped!")
		os.Exit(1)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", "", "Port to listen on")
	rootCmd.Flags().StringVarP(&origin, "origin", "o", "", "Origin server to proxy requests to")

	rootCmd.MarkFlagsRequiredTogether("origin", "port")
}
