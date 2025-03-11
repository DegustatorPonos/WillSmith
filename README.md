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
- ```:q``` to quit. Ctrl-c also works but does not clear the screen
- ```:r``` to reload the current page
- ```:<number>``` to go to the link by it's index (is placed in the square brackets after every link)

## Collisions

With that navigation scheme some links will colline with common functions. Here are some replacements:
- ```/``` -> ```//```
- ```..``` -> ```../```

# BUILDING
- Specify the home page with HomePage variable in Src/Home.go (temporary solution)
- run ```go run .``` in Src/

# TODO:
- [x] Add a navigation between different pages
- [x] Add scrolling support
- [x] Index page
- [x] Custom error code messages
- [ ] Files downloads
- [ ] Cashing pages
- [ ] Better navigation(Go to topics, etc)
- [ ] More commands to control flow

## Long-term plans:
- [ ] Async page loading
- [ ] Auto-translate some static text/html pages to text/hemeni
