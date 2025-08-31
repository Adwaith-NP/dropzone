# DropZone - File Sharing Software
<img width="697" height="413" alt="Screenshot 2025-08-30 at 4 58 56 PM" src="https://github.com/user-attachments/assets/9f01573a-141b-4525-bb26-79c6dd93ad25" />

DropZone is a cross-platform file sharing software that works on **Linux**, **Windows**, and **Mac**.  
It uses **UDP** for discovering receivers (sharing IP and Drop Name) and **TCP** for transferring files.


## Download
You can download the latest compiled versions from the
[Releases page](https://github.com/adwaith-np/dropzone/releases).

- Linux: `dropzone-linux`
- Mac: `dropzone-mac`
- Windows: `dropzone.exe`

## 🚀 Installation

### Linux (64-bit)
```bash
curl -L https://github.com/adwaith-np/dropzone/releases/download/v1.0.0/dropzone-linux-amd64 -o dropzone-linux-amd64
```
```bash
sudo mv dropzone-linux-amd64 /usr/local/bin/dropzone
```
```bash
chmod +x /usr/local/bin/dropzone
```

### Linux (32-bit)
```bash
curl -L https://github.com/adwaith-np/dropzone/releases/download/v1.0.0/dropzone-linux-386 -o dropzone-linux-386
```
```bash
sudo mv dropzone-linux-386 /usr/local/bin/dropzone
```
```bash
chmod +x /usr/local/bin/dropzone
```

### macOS (Intel, 64-bit)
```bash
curl -L https://github.com/adwaith-np/dropzone/releases/download/v1.0.0/dropzone-mac-amd64 -o dropzone-mac-amd64
```
```bash
sudo mv dropzone-mac-amd64 /usr/local/bin/dropzone
```
```bash
chmod +x /usr/local/bin/dropzone
```

### macOS (Apple Silicon M1/M2/M3)
```bash
curl -L https://github.com/adwaith-np/dropzone/releases/download/v1.0.0/dropzone-mac-arm64 -o dropzone-mac-arm64
```
```bash
sudo mv dropzone-mac-arm64 /usr/local/bin/dropzone
```
```bash
chmod +x /usr/local/bin/dropzone
```

### Windows
Download the .exe file from the Releases page,
then run it directly from PowerShell or Command Prompt:
```bash
.\dropzone-windows-amd64.exe
```
If you want it globally available, rename dropzone-windows-amd64.exe to dropzone.exe and
add the folder containing dropzone.exe to your PATH environment variable.

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

5. Use a custom port:
```bash
dropzone -s -f ./file.txt -p 9090 
dropzone -r -p 9090 
```

6. Use a custom name:
```bash
dropzone -r -n MyDropZone
```

7. Use a custom ip:
```bash
dropzone -s -f ./file.txt -ip 198.162.87.12
dropzone -r -ip 198.162.87.12
```

8. Use a custom download directory:
```bash
dropzone -r -dd ./directory
```



---

## Project Info
- **Project**: DropZone - File Sharing Software
- **Transport**: 
  - **UDP** → Used for sharing IP and DropZone name (discovery).
  - **TCP** → Used for file transfers.

## How DropZone Works

### 1. Network Setup
- Both the **sender** and the **receiver** must be connected to the same network (LAN/Wi-Fi).

### 2. Receiver Broadcast
- The **receiver** broadcasts its DropZone name and IP address over the network.

### 3. Sender Discovery
- The **sender** listens for these broadcasts.  
- A list of available receivers is displayed to the sender.

### 4. Receiver Selection
- The **sender selects** a receiver from the list.

### 5. Metadata Sharing
- The sender sends **file metadata** (file name, size, or directory structure) to the receiver.  
- The sender waits for the receiver to respond.

### 6. Receiver Review
- The **receiver** receives the metadata.  
- Details are displayed:
  - If it’s a **file**: the receiver can view the file name and size (or skip).  
  - If it’s a **directory**: the receiver can view the full directory tree with file sizes (or skip).  

### 7. Acceptance Decision
- The receiver is asked to **accept** or **reject** the transfer.  
- If **rejected** → the process ends.  
- If **accepted** → the download begins.  

## 📜 License

This project is licensed under the **MIT License** – feel free to use, modify, and distribute.
