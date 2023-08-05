export class Board {
  constructor(
    readonly board: BoardInfo[][],
    readonly boardImg: string[]
  ) {}
}

export class Koma {
  constructor(
    readonly turn: number,
    readonly x: number,
    readonly y: number,
    readonly size: number,
  ) {}
}

export class BoardInfo {
  size: number;
  turn: number;
  constructor(
    turn: number,
    size: number,
  ) {
    this.size = size;
    this.turn = turn;
  }
}