import { Board, Koma } from "../domain/entity";
import { GameOutPutPort } from "./port/outputPort";

export class GameUseCase {
    constructor(readonly gameoutputport: GameOutPutPort){}

    input(index: number,turn: number, b: number[][],bi: string[]) {
        const res = this.convert(index);
        const koma = new Koma(turn,res.x,res.y)
        const board = new Board(b,bi);

        if(!this.isEmpty(board,koma)){
            return new Koma(-1,-1,-1);    
        }

        b[koma.x][koma.y] = turn;
        if(turn === 1) {
            bi[index] = 'red.200';
        } else if(turn === 2) {
            bi[index] = 'blue.200';
        }
        const newboard = new Board(b,bi);
        this.gameoutputport.display(newboard);

        return new Koma(1,res.x,res.y);
    }

    init() {
        const b: number[][] = [] as number[][]; 
        const bi:string[] = [] as string[];
        
        for (let i = 0; i < 5; i++) {
            b.push([]);
            for (let j = 0; j < 5; j++) {
                b[i].push(0);
            }
        }

        for(let i = 0; i < 25; i++){
            bi[i] = '';
        }
        
        const board = new Board(b,bi);
        this.gameoutputport.display(board);
    }

    checkVertical(board: Board) {
        for(let i = 0; i < 5; i++) {
            if(board.board[0][i] == board.board[1][i] && board.board[1][i] == board.board[2][i] &&
               board.board[2][i] == board.board[3][i] && board.board[3][i] == board.board[4][i] &&
               board.board[0][i] != 0) {
                return true;
            }
        }
        return false;
    }

    checkHorizon(board: Board) {
        for(let i = 0; i < 5; i++) {
            if(board.board[i][0] == board.board[i][1] && board.board[i][1] == board.board[i][2] &&
               board.board[i][2] == board.board[i][3] && board.board[i][3] == board.board[i][4] &&
               board.board[i][0] != 0) {
                return true;
            }
        }
        return false;
    }

    checkCross(board: Board) {
        if( board.board[0][0] == board.board[1][1] && board.board[1][1] == board.board[2][2] &&
            board.board[2][2] == board.board[3][3] && board.board[3][3] == board.board[4][4] &&
            board.board[0][0] != 0) {
            return true;
        }

        return false;
    }

    isWin(inputBoard: number[][]) {
        const bi: string[] = [];
        const board = new Board(inputBoard,bi);
        if(this.checkVertical(board) || this.checkHorizon(board) || this.checkCross(board)) {
            return true;
        }
        return false;
    }

    isEmpty(board: Board, koma: Koma) {
        if(board.board[koma.x][koma.y] != 0) {
            return false;
        }
        return true;
    }

    convert(index: number):InputData {
        if(index >= 0 && index <= 4) {
            return new InputData(0,index);
        } else if(index >= 5 && index <= 9 ) {
            return new InputData(1,index - 5);
        } else if(index >= 10 && index <= 14) {
            return new InputData(2,index - 10);
        } else if(index >= 15 && index <= 19) {
            return new InputData(3,index - 15);
        } else if(index >= 20 && index <= 24) {
            return new InputData(4,index - 20);
        } else {
            return new InputData(-1,-1);
        }
    }
}

class InputData {
    constructor(
        readonly x: number,
        readonly y: number,
    ){}
}