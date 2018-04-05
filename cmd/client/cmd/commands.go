// Copyright © 2018 Florian Lloyd-Poetscher
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolP("force", "f", false, "Forces deletion of all redirects for a hostname")
}

var pingCmd = &cobra.Command{
	Use:     "ping",
	Aliases: []string{"test"},
	Short:   "ping the server to test connection",
	Args:    cobra.NoArgs,
	Hidden:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestFromServer("ping", args)
	},
}

var listCmd = &cobra.Command{
	Use:     "list [hostname] [url]",
	Aliases: []string{"show"},
	Short:   "list the redirects on the server",
	Long: `list shows the redirects on a server. 

	The command allows three forms
	  list                Lists all redirects
	  list hostname       Lists all redirects for a specific hostname
	  list hostname url   Shows the redirect for a specific hostname and url
	`,
	Example: "list www.example.com /",
	Args:    cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestFromServer("list", args)
	},
}

var addCmd = &cobra.Command{
	Use:     "add hostname url target",
	Aliases: []string{"insert", "change"},
	Short:   "adds a new or changes an existing redirect",
	Example: "add www.example.com / http://www.google.com",
	Args:    cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestFromServer("add", args)
	},
}

var removeCmd = &cobra.Command{
	Use:     "remove hostname [url]",
	Aliases: []string{"delete"},
	Short:   "remove a redirect from the server",
	Long: `The command allows two forms
	  remove hostname       Removes all redirects for a hostname 
	  remove hostname url   Removes the redirect for a specific hostname and url
	`,
	Example: "remove www.example.com /",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			// full host should be removed
			hostname := args[0]
			forced, err := cmd.Flags().GetBool("force")
			if err != nil {
				return err
			}
			if !forced {
				fmt.Printf("Confirm to delete all redirects from hostname %v (y/n) ", hostname)
				var answer string
				_, err = fmt.Scanf("%s.1", &answer)
				if err != nil {
					return fmt.Errorf("Could not read from keyboard")
				}
				if answer != "y" && answer != "yes" {
					fmt.Printf("Deletion aborted for hostname %v \n", hostname)
					return nil
				}
			}
			return requestFromServer("deleteHost", args)
		}
		return requestFromServer("delete", args)
	},
}
