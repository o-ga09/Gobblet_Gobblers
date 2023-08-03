import { when } from 'jest-when';
import { Board, Koma } from '../../domain/entity';
import { GameUseCase } from '../gemausecase';
import { GameOutPutPort } from '../port/outputPort';

describe('ユースケースのテスト', () => {
  test('入力を受け取る', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg1 = [
      [0, 0, 0],
      [0, 0, 0],
      [0, 0, 0]
    ];
    const arg2 = ['white', 'white', 'white', 'white', 'white', 'white', 'white', 'white', 'white'];
    const actual = usecase.input(1, 1, arg1, arg2);
    const expected = new Koma(1, 0, 1);
    expect(actual).toEqual(expected);
  });

  test('盤面の配列情報を初期化する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);

    const usecase = new GameUseCase(outputport);
    usecase.init();
    expect(displayMock).toBeCalledTimes(1);
  });

  test('任意の縦一列が揃ったかを判定する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg1 = [
      [1, 0, 0],
      [1, 0, 0],
      [1, 0, 0]
    ];
    const arg2 = ['white', 'white', 'white', 'white', 'white', 'white', 'white', 'white', 'white'];
    const board = new Board(arg1, arg2);
    const actual = usecase.checkVertical(board);
    const expected = true;
    expect(actual).toEqual(expected);
  });

  test('任意の横一列が揃ったかを判定する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg1 = [
      [1, 1, 1],
      [0, 0, 0],
      [0, 0, 0]
    ];
    const arg2 = [
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white',
      'white'
    ];
    const board = new Board(arg1, arg2);
    const actual = usecase.checkHorizon(board);
    const expected = true;
    expect(actual).toEqual(expected);
  });

  test('任意の斜め一列が揃ったかを判定する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg1 = [
      [1, 0, 0, 0, 0],
      [0, 1, 0, 0, 0],
      [0, 0, 1, 0, 0],
      [0, 0, 0, 1, 0],
      [0, 0, 0, 0, 1]
    ];
    const arg2 = ['white', 'white', 'white', 'white', 'white', 'white', 'white', 'white', 'white'];
    const board = new Board(arg1, arg2);
    const actual = usecase.checkCross(board);
    const expected = true;
    expect(actual).toEqual(expected);
  });

  test('任意の一列が揃っていた場合、終了処理を実行する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg = [
      [1, 0, 0],
      [1, 0, 0],
      [1, 0, 0]
    ];
    const actual = usecase.isWin(arg);
    const expected = true;
    expect(actual).toEqual(expected);
  });

  test('駒を置けるかを判定する', () => {
    const outputport = {} as GameOutPutPort;
    const displayMock = jest.fn();
    outputport.display = displayMock;
    when(displayMock).calledWith().mockReturnValueOnce(null);
    const usecase = new GameUseCase(outputport);

    const arg1 = [
      [1, 1, 1],
      [1, 1, 1],
      [1, 1, 1]
    ];
    const arg2 = ['white', 'white', 'white', 'white', 'white', 'white', 'white', 'white', 'white'];
    const board = new Board(arg1, arg2);
    const koma = new Koma(1, 1, 1);
    const actual = usecase.isEmpty(board, koma);
    const expected = false;
    expect(actual).toEqual(expected);
  });
});
