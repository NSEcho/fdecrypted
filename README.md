# fdecrypted
Download file or executable from the application bundle directory using frida-go.

# Installation
```bash
$ go install github.com/lateralusd/fdecrypted@latest
```

# Usage
Filename is relative to the application bundle directory.

## file mode

`fdecrypted APP:DIR:FILENAME`

Possible `DIR` flags include:
* `B` - Applications bundle path
* `D` - Applications directory path
* `L` - Applications library path

Example downloading `some_filename.json` from the application bundle that the Gadget is attached to.

```bash
$ fdecrypted Gadget:B:some_filename.json
```

Example downloading plist file from library directory:

```bash
$ fdecrypted Gadget:L:Preferences/com.example.app.plist
```

## executable mode

`fdecrypted APP`

Example download binary that the gadget is attached to.

```bash
$ fdecrypted Gadget
```
