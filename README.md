# DropZone - File Sharing Software

DropZone is a cross-platform file sharing software that works on **Linux**, **Windows**, and **Mac**.  
It uses **UDP** for discovering receivers (sharing IP and Drop Name) and **TCP** for transferring files.

---

## Features
- Cross-platform support
- Fast file transfers using TCP
- Auto-discovery of receivers using UDP
- Flexible options for file and directory transfers

---

## Usage

### Command Line Options

```
  -d string
        Directory path
  -dd string
        Custom download directory (default "/Users/adwaith/Downloads")
  -f string
        File path
  -fl string
        Comma-separated list of file paths
  -h    Show help
  -ip string
        Custom IPV4 setup (default "192.168.31.48")
  -n string
        Custom drop name setup (default "DropZone")
  -p int
        Custom port setup <warning: also set same port on receiver or sender> (default 8080)
  -r    Run in receiver mode
  -s    Run in sender mode
  -version
        Show version
```

---

## Example Commands

1. Run in **Receiver Mode**:
```bash
dropzone -r
```

2. Run in **Sender Mode** with a single file:
```bash
dropzone -s -f ./file.txt
```

3. Send a whole directory:
```bash
dropzone -s -d ./myfolder
```

4. Send file list:
```bash
dropzone -s -fl ./myfolder,./myfolder2,./myfolder3
```

4. Use a custom port:
```bash
dropzone -s -f ./file.txt -p 9090 
```

5. Use a custom name:
```bash
dropzone -s -f ./file.txt -n MyDropZone
```

5. Use a custom ip:
```bash
dropzone -s -f ./file.txt -ip 198.162.87.12
```

---

## Project Info
- **Project**: DropZone - File Sharing Software
- **Transport**: 
  - UDP → Share IP and Drop Name (discovery)
  - TCP → Share files (transfer)