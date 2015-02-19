# lubeck

## Setup
You'll need go setup with cross compiling
  
  brew install go --with-cc-all

## Build
  
  GOARM=6 GOARCH=arm GOOS=linux go build main.go
