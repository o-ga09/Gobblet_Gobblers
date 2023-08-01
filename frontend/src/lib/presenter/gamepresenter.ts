import { GameState } from "../../view/state/gamestate";
import { Board } from "../domain/entity";
import { GameOutPutPort } from "../usecase/port/outputPort";

export class GamePresenter implements GameOutPutPort {
    constructor(readonly state: GameState){}
    
    display(board: Board): void {
        this.state.setBoard(board);
    }
}