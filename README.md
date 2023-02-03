# fdecrypted
Download file or executable from the application bundle directory using frida-go.

# Installation
```bash
$ go install github.com/lateralusd/fdecrypted
```

# Usage
Filename is relative to the application bundle directory.

## file mode

`fdecrypted APP:FILENAME`

Example downloading `some_filename.json` from the application that the Gadget is attached to.

```bash
$ fdecrypted Gadget:some_filename.json
```

## executable mode

`fdecrypted APP`

Example download binary that the gadget is attached to.

```bash
$ fdecrypted Gadget
```