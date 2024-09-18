# Image to Executable

## Quick start

### If you need to make an executable
```bash
go mod verify && go mod tidy && go build
```

### When you have the executable
```
imageToExecutable <inputfile> <outputdir>
```

This will create the directory you specified and write 3 Go files in there:
```
drawImage.go
params.go
main.go
```

If you just want to share that with people that have the same CPU/architecture and OS as you have, then just go ahead and 

```
go build
```

Which will give you an exectuable called after the directory you specified with ```outputdir``` (Windows executables have a .exe suffix)

If you want to share with others... I am afraid you need to learn about [cross compiling](https://en.wikipedia.org/wiki/Cross_compiler) and more specifically [cross compiling with Go](https://opensource.com/article/21/1/go-cross-compiling) or you can just do a trial & error using one or more of the below lines:

```bash
env GOOS=darwin GOARCH=amd64 go build -a -o darwin_amd64_image_to_binary.bin
env GOOS=darwin GOARCH=arm64 go build -a -o darwin_arm64_image_to_binary.bin

env GOOS=windows GOARCH=amd64 go build -a -o windows_amd64_image_to_binary.exe
env GOOS=windows GOARCH=arm64 go build -a -o windows_arm64_image_to_binary.exe
env GOOS=windows GOARCH=arm go build -a -o windows_arm_image_to_binary.exe
env GOOS=windows GOARCH=386 go build -a -o windows_386_image_to_binary.exe

env GOOS=linux GOARCH=386 go build -a -o linux_386_image_to_binary.bin
env GOOS=linux GOARCH=arm go build -a -o linux_arm_image_to_binary.bin
env GOOS=linux GOARCH=arm64 go build -a -o linux_arm64_image_to_binary.bin
env GOOS=linux GOARCH=amd64 go build -a -o linux_amd64_image_to_binary.bin
```


## About this method

This method deconstructs an image and then uses the RGB pixel data to create a source file that has the draw methods to reproduce the image. It uses the template files in the template dir to create a full set of source files.
This source file is then used to compile into a stand alone executable. When this executable is run, it will display the image.

### Advantages

Since this produces an executable, a scanner for images will not check this. The binary must be run to view the image.
Rendering is very fast.
The source code can be sent instead of the binary.
Furthermore, a recipient does not need any special software to view it.

### Disadvantages

The resulting file is large and it is an executable. Directly transferring executables is severely discouraged (and rightly so) by most modern Operating Systems.
A transfer usually will involve compressing the executable in a (password protected) archive file (a ZIP file).

As an executable, it is detectable by the programming language (Go) and the way this language produces binaries. Binaries are also highly detectable via hashes (locality hashes will almost certainly detect this, even if some pixels are changed)

If the recipient does not know the sender, then malware like code can be implemented and running the executable will compromise your system.