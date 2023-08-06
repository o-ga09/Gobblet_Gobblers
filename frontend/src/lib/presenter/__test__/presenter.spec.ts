import { when } from 'jest-when';
import { GameState } from '../../../view/state/gamestate';
import { Board, BoardInfo } from '../../domain/entity';
import { GamePresenter } from '../gamepresenter';

describe('プレゼンターのテスト', () => {
  test('盤面を表示する', () => {
    const state = {} as GameState;
    const displayMock = jest.fn();
    state.setBoard = displayMock;

    const arg = [
      [new BoardInfo(0,0),new BoardInfo(0,0),new BoardInfo(0,0)],
      [new BoardInfo(0,0),new BoardInfo(0,0),new BoardInfo(0,0)],
      [new BoardInfo(0,0),new BoardInfo(0,0),new BoardInfo(0,0)]
    ];
    const arg2 = ['white', 'white', 'white', 'white', 'white', 'white', 'white', 'white', 'white'];
    const arg3 = ['', '', '', '', '', '', '', '', ''];
    const board = new Board(arg, arg2,arg3);
    when(displayMock).calledWith(board);

    const presenter = new GamePresenter(state);
    presenter.display(board);
    expect(displayMock).toBeCalledTimes(1);
    expect(displayMock).toBeCalledWith(board);
  });
});
