docker run --rm -it -v /Users/mq.tran/Development/go:/go -w /go/src/pmd/brightcove golang:1.7.1-onbuild sh -c '
for GOOS in darwin linux windows; do
#  for GOARCH in 386 amd64; do
  for GOARCH in amd64; do
    echo "Building $GOOS-$GOARCH"
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    go build -o bin/brightcove-$GOOS-$GOARCH
  done
done
'
