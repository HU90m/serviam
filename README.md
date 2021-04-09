# Serviam

This is a collection of scripts for getting information about film and show files from TMDB,
as well as a server for browsing and watching. These notes are made for a future me so are unlikely to be useful for anyone else. If you have any questions, please contact me.

## Build

To build install the standard go toolchain and run `go build`.

The `posterplucker` and `getshow` script expect the sxiv image viewer program
to be installed on the system.

## Issues

Currently `os.Rename` is used for moving files which just changes the hardlink and can't copy files from one partition to another. I decided not to fix this and just copy the directories I am working on to the same partition.

## The Server

The media directory should be linked to ./media in serviam's the directory,
then run the following.

```bash
./serviam
```

## Scripts

### posterplucker

This script searches TMDB for the name of the files given,
creates a json file of the same name containing film information
and puts pictures associated with the films
in the directory called `pictures`.
If the film is a member of a collection,
it information about the collection and pictures of it will be downloaded
and put in a `collections` directory.

```bash
posterplucker/posterplucker $(cat misc/api_key) *.{mp4,mkv}
```

where `*.{mp4,mkv}` finds returns the film files
and the apikey is stored in `misc/api_key`.


Usually neater to run from inside film folder as shown below.
so output is contained.

```bash
pushd /path/to/films
~1/posterplucker/posterplucker $(cat ~1/misc/api_key) *.{mp4,mkv}
```

Then follow the script prompts. The output is shown below.

```
2021/04/08 12:52:57 Current query is 'Chicken.Run.2000'.
Are you happy with this query? (y/n)
n
2021/04/08 12:53:01 Type your new query and press enter:
Chicken Run
2021/04/08 12:53:12 Current query is 'Chicken%20Run'.
Are you happy with this query? (y/n)
y
2021/04/08 12:53:14 Searching the TMDB data base for 'Chicken%20Run'.
2021/04/08 12:53:15 There are 3 results for the query 'Chicken%20Run'.
  2: Poultry in Motion: The Making of 'Chicken Run' (2000-06-24)
  1: Chicken Run 2 ()
  0: Chicken Run (2000-06-21)
Can you see the film you want? (y/n)
y
Select a film:
0
2021/04/08 12:53:22 The film 'Chicken Run' has been selected.
2021/04/08 12:53:22 Downloading 'pictures/8XgmIsbpVamdhwcNVsYzohFZTtT.jpg'.
2021/04/08 12:53:22 Downloaded 'pictures/8XgmIsbpVamdhwcNVsYzohFZTtT.jpg' of size 489227.
2021/04/08 12:53:22 Displaying 'pictures/8XgmIsbpVamdhwcNVsYzohFZTtT.jpg'.
Do you confirm this is the correct film? (y/n)
y
2021/04/08 12:55:05 Downloading 'pictures/e8fwYYDYmcht3YHIK6eTO4VfYPU.jpg'.
2021/04/08 12:55:06 Downloaded 'pictures/e8fwYYDYmcht3YHIK6eTO4VfYPU.jpg' of size 91041.
2021/04/08 12:55:06 Getting Film Data.
2021/04/08 12:55:06 Saving Film Data.
2021/04/08 12:55:06 Saving the 'ChickenRunCollection' collection.
2021/04/08 12:55:06 Current query is 'Fear.And.Loathing.In.Las.Vegas.1998'.
Are you happy with this query? (y/n)
y
2021/04/08 12:55:17 Searching the TMDB data base for 'Fear.And.Loathing.In.Las.Vegas.1998'.
2021/04/08 12:55:18 There are 0 results for the query 'Fear.And.Loathing.In.Las.Vegas.1998'.
Couldn't find anything for this query.
Would you like to give up? (y/n)
n
2021/04/08 12:55:24 Type your new query and press enter:
Fear and Loathing
2021/04/08 12:55:39 Current query is 'Fear%20and%20Loathing'.
Are you happy with this query? (y/n)
y
2021/04/08 12:55:42 Searching the TMDB data base for 'Fear%20and%20Loathing'.
2021/04/08 12:55:42 There are 10 results for the query 'Fear%20and%20Loathing'.
  9: Fear and Loathing in Inner Mongolia (2015-01-01)
  8: Fear and Loathing and Party in Las Ponta Delgada ()
  7: Free Lisl: Fear & Loathing in Denver (2006-11-17)
  6: Spotlight on Location: Fear and Loathing in Las Vegas (1998-11-17)
  5: Fear and Loathing in Materndorf (2004-10-26)
  4: ICW Fear and Loathing XI (2018-12-02)
  3: Fear and Loathing on the Road to Hollywood (1978-11-02)
  2: ICW Fear and Loathing IX (2016-11-20)
  1: Fear & Loathing (2018-07-10)
  0: Fear and Loathing in Las Vegas (1998-05-22)
Can you see the film you want? (y/n)
y
Select a film:
0
2021/04/08 12:55:48 The film 'Fear and Loathing in Las Vegas' has been selected.
2021/04/08 12:55:48 Downloading 'pictures/jwUnGcLBzKNEIzgUUVWAUSwuuBt.jpg'.
2021/04/08 12:55:48 Downloaded 'pictures/jwUnGcLBzKNEIzgUUVWAUSwuuBt.jpg' of size 98591.
2021/04/08 12:55:48 Displaying 'pictures/jwUnGcLBzKNEIzgUUVWAUSwuuBt.jpg'.
Do you confirm this is the correct film? (y/n)
y
```


### posterplacer

This script takes the files json files created by the posterplucker script
and moves the associated pictures and film files of the same name
to the media directory. If a film is a member of a collection,
it will make the collection if it doesn't already exist
and add the film to the collection.

The script expects to be run in the same directory
as the posterplucker script;
It looks for local `picture` and `collection` directories in the same folder.

```bash
~1/posterplacer/posterplacer ../Videos/media/ *.json
```

All moved film's old data files will be placed into a directory called `moved`.

### getshow

Is essentially a combined posterplucker and posterplacer script for shows.

```bash
~1/getshow/getshow $(cat ~1/misc/api_key) Videos/media/ TheBoys/
```
