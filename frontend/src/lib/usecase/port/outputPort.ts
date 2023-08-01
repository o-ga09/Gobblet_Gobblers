import { Board } from "../../domain/entity";

export interface GameOutPutPort {
    display(board: Board): void;
}