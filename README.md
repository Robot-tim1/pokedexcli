## Pokedexcli

A pokedex cli made with the [Pokeapi](https://pokeapi.co/) for [Boot.dev](https://www.boot.dev) 

## Unique Features 

Most of what's here is made with the guidance of the Boot.dev course, though most implementation details were handled by me. The course told me what to add and sometimes gave a little help to get me started 

([my boot.dev profile here](https://www.boot.dev/u/robottim)) - I'm a free user so don't expect to see much there 

But let me tell you what I've added all on my own past the curriculum

### Command Line Shortcuts:

With no libraries I mapped and handled standard Emacs/Bash keyboard shortcuts manually using the raw byte streams from key inputs when the terminal is set to raw mode

Here is the current list of terminal features
- Cycling through your command history with up and down arrow keys
- Navigating the cursor with left and right arrow keys
- Handles keyboard inputs as expected no matter where the cursor is
- Moving the cursor to the ends with Ctrl+A and Ctrl+E
- Clearing the screen with Ctrl+L
- Deleting on either ends of the cursor with Ctrl+U and Ctrl+K
- Pasting what you killed with Ctrl+U/K with Ctrl+Y

### Persistent Save data:

Every time you exit, your pokedex is saved as a json file inside a pokedexcli folder which is created at documents or .local/share depending on OS

When the program starts again it reads from the file and updates the pokedex to have the save data, which can be deleted with the delete command

## Commands

| Command | Description |
|---------|-------------|
| `help` | Displays a help message |
| `exit` | Exit the Pokedex |
| `map` | Displays next locations |
| `mapb` | Displays previous locations |
| `explore <location>` | Type a location name after explore. Displays pokemon in an area |
| `inspect <pokemon>` | Enter a pokemon's name after the command. Displays data about pokemon |
| `catch <pokemon>` | Enter the pokemon's name after the command. Allows you to attempt to catch a pokemon |
| `pokedex` | Displays all pokemon you have caught |
| `delete` | Deletes save data |

## Installation

This is only meant to be a showcase where you can look at my code and judge me

I don't plan on releasing this, but you can still clone the project if you so please
