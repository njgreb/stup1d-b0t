# stup1d-b0t
stup1d-b0t is a discord bot that does stuff.

Current stuffs:
 1. Tells you weather (via open weather api)
 2. Shows you stupid gifs (via tenor)
 3. Lets you rename it (this is awful...really...anyone????)

Planned stuffs:
 1. Channel admin
 2. Cowsay
 3. Google search
 4. Hearthstone stuffs (card lookup, decks, etc?)
 5. Rock/Paper/Scissors game
 6. Youtube search
 7. Probably some other stuff...

## This whole thing is a big ol' work in progress
So yeah, at the moment I am just screwing around with this and Golang.  So, take that in mind.  Happy to accept help from others.

## Required ENV vars
| Key | Value |
|--|--|
| redisHost | hostname or ip address of your redis server |
| redisPort | port redis is running on |
| botToken | Valid Discord bot token |
| weatherApiToken | Valid open weather API token |
| commandPrefix | Prefix for all commands (example, . makes .w work) |
| tenorKey | Valid tenor API key |
