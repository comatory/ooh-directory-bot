#!/bin/sh
STAGED_GO_FILES=$(git diff --cached --name-only | grep ".go$")

PASS=true

for FILE in $STAGED_GO_FILES
do
    go fmt $FILE
done

if ! PASS; then
    printf "COMMIT FAILED\n"
    exit 1
fi

exit 0
