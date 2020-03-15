package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	_extensionHomeURL  = "https://extensions.gnome.org"
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
	version                 print gnomex version
	upgrade [uuid]...       upgrade extension
	about <uuid>            print detailed information of the extension
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

	Upgrade all extensions
	$ gnomex upgrade

	Upgrade some extensions
	$ gnomex dash-to-dock@micxgx.gmail.com Resource_Monitor@Ory0n

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

func checkArgs(badCondition bool) {
	if badCondition {
		fmt.Println("unknown arguments")
		fmt.Println("type `gnomex help` to see usage")
		os.Exit(1)
	}
}

func (g *gnomex) run() {
	if len(os.Args) == 1 {
		fmt.Print(_helpText)
		return
	}

	command := os.Args[1]

	switch command {
	case "search":
		checkArgs(len(os.Args) > 3)

		query := ""
		if len(os.Args) == 3 {
			query = os.Args[2]
		}

		g.search(query)
	case "list":
		checkArgs(len(os.Args) != 2)
		g.list()
	case "install":
		checkArgs(len(os.Args) != 3)
		g.install(os.Args[2])
	case "uninstall":
		checkArgs(len(os.Args) != 3)
		g.uninstall(os.Args[2])
	case "upgrade":
		if len(os.Args) == 2 {
			// upgradeAll()
		} else if len(os.Args) > 2 {
			for _, UUID := range os.Args[2:] {
				_ = UUID
				// upgrade(UUID)
			}
		}
	case "about":
		checkArgs(len(os.Args) != 3)
		g.about(os.Args[2])
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

func printShortInfo(v Extension) {
	color.Yellow.Print(v.Name)
	color.Green.Print(" (" + v.UUID + ") ")
	color.Magenta.Print("by ")
	color.Cyan.Println(v.Creator)
}

func (g *gnomex) search(query string) {
	g.fetchDb(query)
	for _, v := range g.extensions {
		printShortInfo(v)
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

	printShortInfo(extn)

	fmt.Println("downloading extension")
	fileName := g.download(UUID)
	defer os.Remove(fileName)

	_, err := exec.Command("gnome-extensions", "install", "--force", fileName).Output()
	if err != nil {
		fmt.Println("\nunable to install extension")
		os.Exit(1)
	}
	fmt.Println("\nextension installed")

	_, err = exec.Command("gnome-extensions", "enable", UUID).Output()
	if err != nil {
		fmt.Println("\nunable to enable extension")
		os.Exit(1)
	}
	fmt.Println("extension enabled")

	fmt.Print("to activate the extension restart GNOME Shell by pressing ")
	color.Yellow.Print("Alt + F2")
	fmt.Print(" and enter ")
	color.Yellow.Println("r")
}

type writeCount int

func (wc *writeCount) Write(p []byte) (int, error) {
	n := len(p)
	*wc += writeCount(n)
	fmt.Printf("\r%s", "                                          ")
	size := (float32(*wc) / 1024) / 1024
	fmt.Printf("\r%.2f MB downloaded", size)
	return n, nil
}

// download downloads the extension with given UUID.
// Returns the location of the downloaded file.
func (g *gnomex) download(UUID string) string {
	extn := g.extensions[UUID]

	downloadURL := strings.Replace(_downloadURLFormat, "UUID", extn.UUID, 1)
	downloadURL = strings.Replace(downloadURL, "@", "", 1)
	version := strconv.Itoa(extn.ShellVersion[g.gnomeShellVersion].Version)
	downloadURL = strings.Replace(downloadURL, "VERSION", version, 1)

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		fmt.Println("unable to form request to search:", err)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:74.0) Gecko/20100101 Firefox/74.0")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	res, err := g.client.Do(req)
	if err != nil {
		fmt.Println("unable to download extension from", downloadURL)
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	tmpFile, err := ioutil.TempFile("", UUID)
	if err != nil {
		fmt.Println("unable to create file to save:", err)
		os.Exit(1)
	}
	defer tmpFile.Close()
	fileName := tmpFile.Name()

	count := writeCount(0)
	if _, err = io.Copy(tmpFile, io.TeeReader(res.Body, &count)); err != nil {
		fmt.Println("unable to write to file:", err)
		os.Exit(1)
	}

	return fileName
}

// unisntall uinstalls the extension with given UUID
func (g *gnomex) uninstall(UUID string) {
	fmt.Println("uinstalling", UUID)
}

func (g *gnomex) upgradeAll() {

}

func (g *gnomex) upgrade(UUID string) {

}

func (g *gnomex) about(UUID string) {
	g.fetchDb(UUID)

	v, ok := g.extensions[UUID]
	if !ok {
		fmt.Println("extension with UUID", UUID, "not found")
		os.Exit(1)
	}

	printShortInfo(v)
	fmt.Printf("%v\n\n%v\n", _extensionHomeURL+v.Link, v.Description)
}
