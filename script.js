rpc.exports = {
    download_file(filename) {
        var bd = ObjC.classes.NSBundle.mainBundle().bundleURL().toString().slice(7);
        bd += filename;
        var dt = ObjC.classes.NSData.alloc().initWithContentsOfFile_(bd);
        var arr = Memory.readByteArray(dt.bytes(), dt.length());
        send(filename, arr);
    },
    download_bin() {
        var execPath = ObjC.classes.NSBundle.mainBundle().executablePath();
        var dt = ObjC.classes.NSData.alloc().initWithContentsOfFile_(execPath);
        var arr = Memory.readByteArray(dt.bytes(), dt.length());
        send(execPath.toString(), arr);
    },
}