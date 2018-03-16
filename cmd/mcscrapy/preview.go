package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	address    = "127.0.0.1"
	port       = "8000"
	previewCmd = &cobra.Command{
		Use:     "preview",
		Aliases: []string{"p"},
		Short:   "Preview the scraped website.",
		Long:    "Run a local webserver to preview a scraped website.",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
			} else {
				runPreview(args[0])
			}
		},
	}
)

func init() {
	previewCmd.PersistentFlags().StringVarP(&address, "address", "a", "127.0.0.1", "")
	previewCmd.PersistentFlags().StringVarP(&port, "port", "p", "8000", "")
}

func runPreview(site string) {
	fmt.Println("Previewing "+site+" at", "http://"+address+":"+port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if filepath.Ext(r.URL.Path) == "" && r.URL.Path != "/" {
			http.ServeFile(w, r, "site/"+site+"/"+strings.TrimRight(r.URL.Path, "/")+".html")
		} else {
			http.ServeFile(w, r, "site/"+site+"/"+r.URL.Path)
		}
	})

	log.Fatal(http.ListenAndServe(address+":"+port, nil))
}
