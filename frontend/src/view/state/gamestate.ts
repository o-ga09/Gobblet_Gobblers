import { Board } from "../../lib/domain/entity";

export class GameState {
    board: number[][] = [];
    boardImg: string[] = [];

    setBoard(b: Board) {
        this.board = b.board;
        this.boardImg = b.boardImg;
    }
}