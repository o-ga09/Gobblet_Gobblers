import { Board, Koma } from '../domain/entity';
import { GameOutPutPort } from './port/outputPort';

const BOARD_COLUMN = 3;
const BOARD_ROW = 3;
const BOARD = 9;
export class GameUseCase {
  constructor(readonly gameoutputport: GameOutPutPort) {}

  input(index: number, turn: number, b: number[][], bi: string[]) {
    const res = this.convert(index);
    const koma = new Koma(turn, res.x, res.y);
    const board = new Board(b, bi);

    if (!this.isEmpty(board, koma)) {
      return new Koma(-1, -1, -1);
    }

    b[koma.x][koma.y] = turn;
    if (turn === 1) {
      bi[index] = 'red.200';
    } else if (turn === 2) {
      bi[index] = 'blue.200';
    }
    const newboard = new Board(b, bi);
    this.gameoutputport.display(newboard);

    return new Koma(1, res.x, res.y);
  }

  init() {
    const b: number[][] = [] as number[][];
    const bi: string[] = [] as string[];

    for (let i = 0; i < BOARD_COLUMN; i++) {
      b.push([]);
      for (let j = 0; j < BOARD_ROW; j++) {
        b[i].push(0);
      }
    }

    for (let i = 0; i < BOARD; i++) {
      bi[i] = '';
    }

    const board = new Board(b, bi);
    this.gameoutputport.display(board);
  }

  checkVertical(board: Board) {
    for (let i = 0; i < BOARD_COLUMN; i++) {
      if (
        board.board[0][i] == board.board[1][i] &&
        board.board[1][i] == board.board[2][i] &&
        board.board[2][i] == board.board[0][i] &&
        board.board[0][i] != 0
      ) {
        return true;
      }
    }
    return false;
  }

  checkHorizon(board: Board) {
    for (let i = 0; i < BOARD_ROW; i++) {
      if (
        board.board[i][0] == board.board[i][1] &&
        board.board[i][1] == board.board[i][2] &&
        board.board[i][2] == board.board[i][0] &&
        board.board[i][0] != 0
      ) {
        return true;
      }
    }
    return false;
  }

  checkCross(board: Board) {
    if (
      board.board[0][0] == board.board[1][1] &&
      board.board[1][1] == board.board[2][2] &&
      board.board[2][2] == board.board[0][0] &&
      board.board[0][0] != 0
    ) {
      return true;
    } else if (
      board.board[0][2] == board.board[1][1] &&
      board.board[1][1] == board.board[2][0] &&
      board.board[0][2] == board.board[2][0] &&
      board.board[0][2] != 0
    ) {
      return true;
    }

    return false;
  }

  isWin(inputBoard: number[][]) {
    const bi: string[] = [];
    const board = new Board(inputBoard, bi);
    if (this.checkVertical(board) || this.checkHorizon(board) || this.checkCross(board)) {
      return true;
    }
    return false;
  }

  isEmpty(board: Board, koma: Koma) {
    if (board.board[koma.x][koma.y] != 0) {
      return false;
    }
    return true;
  }

  convert(index: number): InputData {
    if (index >= 0 && index <= 2) {
      return new InputData(0, index);
    } else if (index >= 3 && index <= 5) {
      return new InputData(1, index - 3);
    } else if (index >= 6 && index <= 8) {
      return new InputData(2, index - 6);
    } else {
      return new InputData(-1, -1);
    }
  }
}

class InputData {
  constructor(
    readonly x: number,
    readonly y: number
  ) {}
}
