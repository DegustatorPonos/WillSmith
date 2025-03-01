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

# BUILDING
- Specify the home page with HomePage variable in Src/Home.go (temporary solution)
- run ```go run .``` in Src/

# TODO:
- [x] Add a navigation between different pages
- [x] Add scrolling support
- [ ] Index page
- [ ] Custom error code messages
- [ ] Cashing pages
