import { when } from 'jest-when';
import { GameState } from '../../../view/state/gamestate';
import { Board } from '../../domain/entity';
import { GamePresenter } from '../gamepresenter';

describe('プレゼンターのテスト', () => {
  test('盤面を表示する', () => {
    const state = {} as GameState;
    const displayMock = jest.fn();
    state.setBoard = displayMock;

    const arg = [
      [0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0],
      [0, 0, 0, 0, 0]
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
    const board = new Board(arg, arg2);
    when(displayMock).calledWith(board);

    const presenter = new GamePresenter(state);
    presenter.display(board);
    expect(displayMock).toBeCalledTimes(1);
    expect(displayMock).toBeCalledWith(board);
  });
});
