export class Board {
    constructor(
        readonly board: number[][],
        readonly boardImg: string[]
    ){}
}

export class Koma {
    constructor(
        readonly turn: number,
        readonly x: number,
        readonly y: number,
    ){}
}