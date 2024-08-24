#!/bin/bash

echo "\n-----ready to run post create commands-----\n"

cd /workspace

# echo "- Removing /dist"
# if [ -d /workspace/dist/ ]; then
#   rm -rf /workspace/dist/;
# fi

elvish

exec "$@"
