# rename all directories in this directory to have a prefix of "go-"
# if they don't already have a prefix of "go-"
# usage: foo.sh <directory> <prefix>
# example: foo.sh go go-

pushd $1
for dir in $(ls); do
    if [ -d "$dir" ]; then
        if [[ "$dir" != $2* ]]; then
            mv "$dir" "$2$dir"
        fi
    fi
done
popd
