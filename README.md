# go_ms #
- - - -
This repo houses the game and solver. The game is CLI-based and the solver logic is described below.

## Solver ##
- - - -
This will house the strategy/logic needed to solve Minesweeper automatically. The current thought process is the following:

1. Identify all safe moves and mines based on immediate neighbors
2. Identify all safe moves based on known mines
3. Identify all safe moves by using neighbors' info
4. Identify least risky move
5. Blind click if no other options are found
