# WillSmith
Custom gemini CLI client for linux. 
## WIP

# Navigation
- ```..``` to go back in history
- Type the link to go to it, e.g. ```=> news/ ``` can be accessed with ```news``` command
- You can type a full link and if it atarts with ```gemeni://``` it will try to access it
- ```/``` scrolls the page 1/2 of the screen up
- ```\``` scrolls the page 1/2 of the screen down (or the first line) 
- ```}``` and ```{``` to hop between the white spaces up and down
- ```]``` and ```[``` to hop between headers up and down
- ```:q``` to quit. Ctrl-c also works but does not clear the screen
- ```:r``` to reload the current page
- ```:r``` to cancel page loading
- ```:<number>``` to go to the link by it's index (is placed in the square brackets after every link)

## Collisions

With that navigation scheme some links will colline with common functions. Here are some replacements:
- ```/``` -> ```//```
- ```..``` -> ```../```

# BUILDING
Currently there is no installation option, but you can run it directly by running ```go run .``` in Src/ directory

# TODO:
## Current major change in progress:
- [ ] Cashing pages

## Before release tasks:
- [x] Add a navigation between different pages
- [x] Add scrolling support
- [x] Index page
- [x] Custom error code messages
- [x] Async page loading
- [ ] Files downloads
- [ ] Better navigation(Go to topics, etc)
- [ ] More commands to control flow

## Long-term plans:
- [ ] Auto-translate some static text/html pages to text/hemeni
