package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	color "gopkg.in/gookit/color.v1"
)

const (
	_searchURL         = "https://extensions.gnome.org/extension-query"
	_downloadURLFormat = "https://extensions.gnome.org/extension-data/UUID.vVERSION.shell-extension.zip"
	_version           = "0.0.1"
	_helpText          = `gnomex version ` + _version + `

Search, install and uinstall GNOME Shell extensions.

Commands
	search [query]          search extensions
	list                    list installed extensions
	install <uuid>          install extension with the uuid
	uinstall <uuid>         uninstall extension with the uuid
	version                 print version
	help                    print this help information

Examples
	Search extension with query "user themes"
	$ gnomex search "user themes"

	Search all extensions
	$ gnomex search

	Install dash-to-dock extension
	$ gnomex install dash-to-dock@micxgx.gmail.com

	Uinstall dash-to-dock extension
	$ gnomex uninstall dash-to-dock@micxgx.gmail.com

	List installed extensions
	$ gnomex list

`
)

// gnomex application
type gnomex struct {
	gnomeShellVersion string
	client            http.Client
	extensions        map[string]Extension
}

func findGnomeShellVersion() string {
	out, err := exec.Command("gnome-shell", "--version").Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Format: GNOME Shell 3.34.3
	v := strings.Replace(string(out), "GNOME Shell", "", 1)
	parts := strings.Split(v, ".")
	return strings.TrimSpace(parts[0] + "." + parts[1])
}

func newGnomex() *gnomex {
	g := &gnomex{
		gnomeShellVersion: findGnomeShellVersion(),
		client: http.Client{
			Timeout: time.Second * 2,
		},
		extensions: make(map[string]Extension),
	}

	return g
}

func (g *gnomex) run() {
	if len(os.Args) == 1 {
		fmt.Print(_helpText)
		return
	}

	command := os.Args[1]

	switch command {
	case "search":
		if len(os.Args) > 3 {
			fmt.Println("unknown arguments")
			fmt.Println("type `gnomex help` to see usage")
			return
		}

		query := ""
		if len(os.Args) == 3 {
			query = os.Args[2]
		}

		g.search(query)
	case "list":
		if len(os.Args) != 2 {
			fmt.Println("unknown arguments")
			fmt.Println("type `gnomex help` to see usage")
			return
		}

		g.list()
	case "install":
		if len(os.Args) != 3 {
			fmt.Println("unknown arguments")
			fmt.Println("type `gnomex help` to see usage")
			return
		}

		g.install(os.Args[2])
	case "uninstall":
		if len(os.Args) != 3 {
			fmt.Println("unknown arguments")
			fmt.Println("type `gnomex help` to see usage")
			return
		}

		g.uninstall(os.Args[2])
	default:
		fmt.Print(_helpText)
	}
}

// fetchDb downloads the details of all extensions from gnome extension website.
// A better approach would be cache the db and provide a command to refresh it if necessary.
func (g *gnomex) fetchDb(query string) {
	page := 0

	for {
		req, err := http.NewRequest("GET", _searchURL, nil)
		if err != nil {
			fmt.Println("unable to form request to search:", err)
			os.Exit(1)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:74.0) Gecko/20100101 Firefox/74.0")

		// params: page=1&shell_version=3.34&search=user%20themes
		q := req.URL.Query()
		page++
		q.Set("page", strconv.Itoa(page))
		q.Set("search", query)
		q.Set("shell_version", g.gnomeShellVersion)
		req.URL.RawQuery = q.Encode()

		res, err := g.client.Do(req)
		if err != nil {
			fmt.Println("unable to search:", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("unable to read search result:", err)
			os.Exit(1)
		}

		v := SearchResult{}
		err = json.Unmarshal(b, &v)
		if err != nil {
			fmt.Println("unable to parse search result:", err)
			fmt.Println(req.URL)
			fmt.Println(string(b))
			os.Exit(1)
		}

		for _, a := range v.Extensions {
			g.extensions[a.UUID] = a
		}

		if v.Numpages == page {
			return
		}
	}
}

func (g *gnomex) search(query string) {
	g.fetchDb(query)
	for _, v := range g.extensions {
		color.Yellow.Print(v.Name)
		color.Green.Print(" (" + v.UUID + ") ")
		color.Magenta.Print("by ")
		color.Cyan.Println(v.Creator)
	}
}

// list lists all installed extensions
func (g *gnomex) list() {
	fmt.Println("listing...")
}

// install installs the extension with given UUID
func (g *gnomex) install(UUID string) {
	g.fetchDb(UUID)
	extn, ok := g.extensions[UUID]
	if !ok {
		fmt.Println("unable to find extension")
		return
	}
	fmt.Println("installing", extn)
}

// download downloads the extension with given UUID
func (g *gnomex) download(UUID string) {
	g.fetchDb(UUID)
	extn, ok := g.extensions[UUID]
	if !ok {
		fmt.Println("unable to find extension")
		return
	}
	fmt.Println("downloading", extn)
}

// unisntall uinstalls the extension with given UUID
func (g *gnomex) uninstall(UUID string) {
	fmt.Println("uinstalling", UUID)
}
