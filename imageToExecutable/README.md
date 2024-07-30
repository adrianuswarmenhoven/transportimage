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