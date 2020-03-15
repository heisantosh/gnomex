# gnomex

A CLI tool to search and install GNOME Shell extensions.

## Context
The usual way to install GNOME Shell extensions is to visit https://extensions.gnome.org/. And install the browser extension and the host connector and install the extension. This tool aims to ease the installation of extensions through a CLI.

### Finding available extensions
A HTTP GET of the URL https://extensions.gnome.org/extension-query/?page=1&shell_version=3.34 returns a list of extensions that support GNOME Shell version 3.34. The response is of the format:

```json
{
    "extensions": [
        {
            "uuid": "user-theme@gnome-shell-extensions.gcampax.github.com",
            "name": "User Themes",
            "creator": "fmuellner",
            "creator_url": "/accounts/profile/fmuellner",
            "pk": 19,
            "description": "Load shell themes from user directory.",
            "link": "/extension/19/user-themes/",
            "icon": "/static/images/plugin.png",
            "screenshot": null,
            "shell_version_map": {
                "3.26": {
                    "pk": 7480,
                    "version": 32
                },
                "3.24": {
                    "pk": 7481,
                    "version": 33
                },
                "3.28": {
                    "pk": 8103,
                    "version": 34
                },
                "3.30": {
                    "pk": 8388,
                    "version": 35
                },
                "3.32": {
                    "pk": 10231,
                    "version": 37
                },
                "3.34": {
                    "pk": 13345,
                    "version": 39
                },
                "3.36": {
                    "pk": 14396,
                    "version": 40
                }
            }
        },
        ...
    ],
    "total": 10,
    "numpages": 31
}
```

To find the extensions by search keyword `user themes` the HTTP GET request is -
https://extensions.gnome.org/extension-query/?page=1&shell_version=3.34&search=user%20themes.

### Finding GNOME Shell version
The current version of GNOME Shell can be found using the command:

```bash
$ gnome-shell --version
GNOME Shell 3.34.1
```

### Version of an extension
An extension can have different versions. Multiple versions of the extension could be use in the same GNOME Shell version.

For example, for the extension `User Themes` by `fmuellner` for GNOME Shell version 3.34 there are 2 versions of the extension - versions 38 and 39.

### Downloading an extension
Below is the HTTP GET URL to downlod the `User Themes` extension -
https://extensions.gnome.org/download-extension/user-theme%40gnome-shell-extensions.gcampax.github.com.shell-extension.zip?version_tag=13345

It is a specific format - https://extensions.gnome.org/download-extension/{`uuid`}.shell-extension.zip?version_tag=`pk`

`version_tag` is same as the field `pk` in the `shell_version_map` field.

Another way to download to GET the URL https://extensions.gnome.org/extension-data/user-themegnome-shell-extensions.gcampax.github.com.v39.shell-extension.zip. 

This URL also is a specific format - https://extensions.gnome.org/extension-data/{`uuid`}.v{`version`}.shell-extension.zip

Here `v39` refers the version of the extension. It's the same as the `version` field in the `shell_version_map` field.

### Installing a downloaded extension

```bash
$ cd ~/Downloads
$ ls
dash-to-dockmicxgx.gmail.com.v67.shell-extension.zip
$ gnome-extensions install dash-to-dockmicxgx.gmail.com.v67.shell-extension.zip
dash-to-dock@micxgx.gmail.com
$ # Enable the extension
$ gnome-extensions enable dash-to-dock@micxgx.gmail.com
```

Then restart GNOME Shell by pression `Alt + F2` and enter `r`.

Now the extension will be active.

### Managing extension settings
`GNOME Tweaks` application can be used to manage the installed extensions and the settings of the extensions.

## Plan

### Searching for extensions
Find the running gnome shell version. Get the results of the query for that version and the search query from the gnome extension website. Store the results in a map of UUID to the extension. Print the list of extensions in the search result.

```bash
$ gnomex search "dock"
name (uuid) by creator
name (uuid) by creator
```

### Installing the extension
Find the extension details from the gnome extension website by querying with the UUID. Get the extension details from the map using UUID as key.

```bash
$ gnomex install uuid
installing extension: done
restarting GNOME Shell: done
extension is ready to use
```

### Uninstalling the extension
```bash
$ gnomex uninstall uuid
uninstalling extension: done
restarting GNOME Shell: done
extension is removed
```

### List installed extensions
```bash
$ gnomex list
uuid - description
uuid - description
uuid - description
```

### Upgrade installed extensions

Upgrade all extensions
```bash
$ gnomex upgrade
```

Upgrade some extensions
```bash
$ gnomex uuid1 uuid2 uuid3
```

### Show detailed information of an extension
```bash
$ gnomex about uuid
```