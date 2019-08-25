#!/usr/bin/bash

echo "{" > index.json
echo "    \"collection\": [" >> index.json

for file in  ~/Videos/YetToWatch/*.mp4
do
    name="$(basename "$file")"
    name_no_ext="$(basename -s .mp4 "$name")"
    u_name="$(basename "${file// /_}")"
    u_name_no_ext="$(basename -s .mp4 "$u_name")"
    
    mkdir "$u_name_no_ext"
    ln -s "$file" "$u_name_no_ext"/"$u_name"

    echo "        {" >> index.json
    echo "            \"id\": \"$u_name_no_ext\"," >> index.json
    echo "            \"name\": \"$name_no_ext\"," >> index.json
    echo "            \"year\": 0," >> index.json
    echo "            \"location\": \"$u_name_no_ext/$u_name\"," >> index.json
    echo "            \"poster\": \"\"" >> index.json
    echo "        }," >> index.json
done
echo "    ]" >> index.json
echo "}" >> index.json
