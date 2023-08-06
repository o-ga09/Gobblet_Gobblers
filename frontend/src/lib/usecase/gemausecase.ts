import { Board, BoardInfo, Koma } from '../domain/entity';
import { GameOutPutPort } from './port/outputPort';

const BOARD_COLUMN = 3;
const BOARD_ROW = 3;
const BOARD = 9;
export class GameUseCase {
  constructor(readonly gameoutputport: GameOutPutPort) {}

  input(index: number, turn: number,size: number, b: BoardInfo[][], bc: string[],bi: string[]) {
    const res = this.convert(index);
    const koma = new Koma(turn, res.x, res.y,size);
    const board = new Board(b, bc,bi);

    if (!this.isEmpty(board, koma)) {
      return new Koma(-1, -1, -1,-1);
    }

    b[koma.x][koma.y].turn = turn;
    b[koma.x][koma.y].size = size;
    if (turn === 1) {
      bc[index] = 'red.200';
    } else if (turn === 2) {
      bc[index] = 'blue.200';
    }

    if(size === 0 || size === 1) bi[index] = '/Gophersvg.svg';             // SMALL
    else if(size === 2 || size === 3) bi[index] = '/Gophersvg_yellow.svg'; // MIDIUM
    else if(size === 4 || size === 5) bi[index] = '/Gophersvg_pink.svg';   // LARGE

    const newboard = new Board(b, bc,bi);
    this.gameoutputport.display(newboard);

    return new Koma(1, res.x, res.y,size);
  }

  init() {
    const b: BoardInfo[][] = [] as BoardInfo[][];
    const bc: string[] = [] as string[];
    const bi: string[]= [] as string[];

    for (let i = 0; i < BOARD_COLUMN; i++) {
      b.push([]);
      for (let j = 0; j < BOARD_ROW; j++) {
        b[i].push( new BoardInfo(0,-1));
      }
    }

    for (let i = 0; i < BOARD; i++) {
      bc[i] = '';
      bi[i] = '';
    }

    const board = new Board(b, bc,bi);
    this.gameoutputport.display(board);
  }

  checkVertical(board: Board) {
    for (let i = 0; i < BOARD_COLUMN; i++) {
      if (
        board.board[0][i].turn == board.board[1][i].turn &&
        board.board[1][i].turn == board.board[2][i].turn &&
        board.board[2][i].turn == board.board[0][i].turn &&
        board.board[0][i].turn != 0
      ) {
        return true;
      }
    }
    return false;
  }

  checkHorizon(board: Board) {
    for (let i = 0; i < BOARD_ROW; i++) {
      if (
        board.board[i][0].turn == board.board[i][1].turn &&
        board.board[i][1].turn == board.board[i][2].turn &&
        board.board[i][2].turn == board.board[i][0].turn &&
        board.board[i][0].turn != 0
      ) {
        return true;
      }
    }
    return false;
  }

  checkCross(board: Board) {
    if (
      board.board[0][0].turn == board.board[1][1].turn &&
      board.board[1][1].turn == board.board[2][2].turn &&
      board.board[2][2].turn == board.board[0][0].turn &&
      board.board[0][0].turn != 0
    ) {
      return true;
    } else if (
      board.board[0][2].turn == board.board[1][1].turn &&
      board.board[1][1].turn == board.board[2][0].turn &&
      board.board[0][2].turn == board.board[2][0].turn &&
      board.board[0][2].turn != 0
    ) {
      return true;
    }

    return false;
  }

  isWin(inputBoard: BoardInfo[][]) {
    const bc: string[] = [];
    const bi: string[] = [];
    const board = new Board(inputBoard, bc,bi);
    if (this.checkVertical(board) || this.checkHorizon(board) || this.checkCross(board)) {
      return true;
    }
    return false;
  }

  isEmpty(board: Board, koma: Koma) {
    if(board.board[koma.x][koma.y].turn !== koma.turn && board.board[koma.x][koma.y].size < koma.size) {
      return true;
    }
    return false;
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
