package main

import "fmt"

const rowLength = 8
const colLength = 8

type game struct {
	board [rowLength][colLength]int
	score map[int]int
	turn  int
}

var theGame = game{
	[8][8]int{},
	map[int]int{},
	1,
}

func setupGame() {
	theGame.board[3][3] = 1
	theGame.board[4][4] = 1
	theGame.board[3][4] = 2
	theGame.board[4][3] = 2
	theGame.score[1] = 2
	theGame.score[2] = 2
	theGame.turn = 1
	printGame(&theGame.board)
	findPotentialMoves(theGame.board, 1)
}

func findPotentialMoves(board [8][8]int, p int) {
	for rowIndex, row := range board {
		for colIndex := range row {
			if board[rowIndex][colIndex] == p {
				checkDirection(0, 1, rowIndex, colIndex, p, &board)
				checkDirection(1, 1, rowIndex, colIndex, p, &board)
				checkDirection(0, -1, rowIndex, colIndex, p, &board)
				checkDirection(1, -1, rowIndex, colIndex, p, &board)
				checkDirection(1, 0, rowIndex, colIndex, p, &board)
				checkDirection(-1, -1, rowIndex, colIndex, p, &board)
				checkDirection(-1, 0, rowIndex, colIndex, p, &board)
				checkDirection(-1, 1, rowIndex, colIndex, p, &board)
			}
		}
	}
	printGame(&board)
}

func checkDirection(offsetY int, offsetX int, originX int, originY int, p int, board *[8][8]int) {
	if moveInBounds(offsetX+originX, offsetY+originY) {
		previousTile := board[originX][originY]
		tile := board[originX+offsetX][originY+offsetY]

		if tile != p && tile != 0 && tile != 3 {
			checkDirection(offsetY, offsetX, originX+offsetX, originY+offsetY, p, board)
		} else if previousTile != p && previousTile != 0 && previousTile != 3 && tile == 0 {
			board[originX+offsetX][originY+offsetY] = 3
		}
	}
}

func movePiece(move moveData) bool {
	var valid = false

	if theGame.board[move.Row][move.Col] == 0 {
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				if i != 0 || j != 0 {
					moveMaid := validateCheckDirection(i, j, move.Row, move.Col, move.Player)
					if moveMaid == true {
						valid = true
					}
				}
			}
		}
	}
	if valid {
		printGame(&theGame.board)
		findPotentialMoves(theGame.board, getOpposingPlayer(move.Player))
		theGame.turn = getOpposingPlayer(theGame.turn)
	}
	return valid
}

func validateCheckDirection(offsetY int, offsetX int, originX int, originY int, p int) bool {
	if moveInBounds(offsetX+originX, offsetY+originY) {
		oP := getOpposingPlayer(p)
		previousTile := theGame.board[originX][originY]
		tile := theGame.board[originX+offsetX][originY+offsetY]

		if tile != p && tile != 0 {
			if validateCheckDirection(offsetY, offsetX, originX+offsetX, originY+offsetY, p) {
				theGame.board[originX][originY] = p
				theGame.score[p] = theGame.score[p] + 1
				theGame.score[oP] = theGame.score[oP] - 1
				return true
			}
		} else if previousTile != p && previousTile != 0 && tile == p {
			theGame.board[originX][originY] = p
			theGame.score[p] = theGame.score[p] + 1
			return true
		}

		return false
	}
	return false
}

func getValueAt(move moveData) {
	if moveInBounds(move.Row, move.Col) {
		fmt.Println(move)
		fmt.Println(theGame.board[move.Row][move.Col])
	} else {
		fmt.Println("Move Out Of Bounds")
	}
}

func getOpposingPlayer(p int) int {
	if p == 1 {
		return 2
	}
	return 1
}

func moveInBounds(x int, y int) bool {
	return x > 0 && y > 0 && x < rowLength-1 && y < colLength-1
}

func printGame(board *[8][8]int) {
	for _, row := range board {
		fmt.Println(row)
	}
	fmt.Printf("Player1: %v - Player2: %v\n", theGame.score[1], theGame.score[2])
}
