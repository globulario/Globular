

# Windows environement setup

How to setup globular developement environnement on windows.

### Install GO
first of all install Go, >the lastest version

https://golang.org/dl/

Test if go is correctly install with 
```bash
go version
> go version go1.20 windows/amd64
```
### Setup Globular code
Now create a dir name globulario
```bash
mkdir globulario
cd globulario
```

Clone git pojet,
```bash
git clone https://github.com/globulario/services.git
git clone https://github.com/globulario/Globular.git
git clone https://github.com/davecourtois/Utility.git
```

Now I will compile the services. The first step is to generate gRpc files.
To do so you must have gRpc installed on your computer. (I will not explain how to install it here).

You will aslo need language specific code generator to generate globular grpc files,

Go
For go simply run,
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
md file not wo
### Typescript

Download protoc-gen-grpc-web from https://github.com/grpc/grpc-web/releases
for the platform you need. Put the exec in the bin path rename it from original name to protoc-gen-grpc-web.exe

Generate gRpc stub...

if you are in a linux shell simply run command 
```bash
sh generateCode.sh
```

> Otherwise copy all the generateCode.sh and paste it in a command prompt, it will run all command one after another.

> C# and C++ generator must be compile from source code, I will not define how to do it but I successfully compile it with msys2 on windows.

Now to compile golang services 
```bash
cd golang
go mod tidy
sh build.sh
```

Now globular from globulario dir,
```bash
cd Globular
go mod tidy
go build
```

Et voila you got a working globular executable and services!

# Linux ARM version (raspberry pie)

Here is a detailed guide on how to set up the development environment for a Raspberry Pi 4.

First, make sure your system is up to date by running the following commands:

```bash
sudo apt update
sudo apt install python3.8
sudo apt install python-is-python3
sudo apt install nano
sudo apt install openssh-server
```

Install the required tools and libraries:

## Git

```bash
sudo apt install git-all
```

### Go > the lastest version

```bash
sudo apt install net-tools
sudo apt-get update
sudo apt-get -y upgrade
mkdir temp
cd temp/
wget https://golang.org/dl/go1.16.7.linux-arm64.tar.gz
sudo tar -xvf go1.16.7.linux-arm64.tar.gz
sudo mv go /usr/local/
sudo nano ~/.bashrc
```

Modify the path in the .bashrc file:
```bash
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
Save the changes and exit the editor.

### SQL unix driver
```bash
sudo apt-get install build-essential
sudo wget http://www.unixodbc.org/unixODBC-2.3.9.tar.gz
sudo tar -xvf unixODBC-2.3.9.tar.gz
cd unixODBC-2.3.9/
sudo ./configure
sudo make all install clean
sudo ldconfig
```

### Protoffuer compiler (protoc)

```bash
sudo apt install -y protobuf-compiler
```

### MongoDB *last version
The rasberry pi version 6.0.5 of mongoDB is not surrported by mongoDB directly but I found
this project who keep binairy availaible,

```bash
wget https://github.com/themattman/mongodb-raspberrypi-binaries/releases/download/r6.0.5-rpi-unofficial/mongodb.ce.pi.r6.0.5.tar.gz 
sudo tar -xvf mongodb.ce.pi.r6.0.5.tar.gz -d /usr/local/bin
```

### Stream Downloader

```bash
sudo apt install python3.8
sudo apt install python-is-python3
# Get the latest binary for your platform from https://github.com/yt-dlp/yt-dlp/releases
sudo wget -qO /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
yt-dlp --version
```
FFmpeg:
Please refer to the link below for instructions on how to install FFmpeg:
https://docs.nvidia.com/video-technologies/video-codec-sdk/ffmpeg-with-nvidia-gpu/

```bash
sudo apt-get install ffmpeg
```

### gRPC code generators
```bash
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### grpcwebproxy
```bash
git clone github.com/improbable-eng/grpc-web.git
cd go/grpcwebproxy
go build
cp grpcwebproxy ~/globulario/services/bin/
sudo chmod +x ~/globulario/services/bin/grpcwebproxy
```

### Globular project itself
```bash
mkdir globulario
cd globulario
git clone https://github.com/globulario/Globular.git
git clone https://github.com/globulario/services.git
git clone https://github.com/davecourtois/Utility.git
cd services
sh generateCode.sh
cd golang
go mod tidy
sh build.sh
cd ..
go mod tidy
go build
```

### Setup environnement variable
In the /etc/environment file, set the ServicesRoot variable to point to the services directory:

```bash
sudo nano /etc/environment
```

set variable, ServicesRoot=/home/ubuntu/globulario/services

Save the changes and exit the editor.
That's it! You have now set up the development environment for a Raspberry Pi 4 and installed the necessary tools and libraries for Globular.
