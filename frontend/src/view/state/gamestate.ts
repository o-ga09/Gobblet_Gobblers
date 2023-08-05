import { Board, BoardInfo } from '../../lib/domain/entity';

export class GameState {
  board: BoardInfo[][] = [];
  boardImg: string[] = [];

  setBoard(b: Board) {
    this.board = b.board;
    this.boardImg = b.boardImg;
  }
}
