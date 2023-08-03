import { GameState } from '../view/state/gamestate';
import { GamePresenter } from './presenter/gamepresenter';
import { GameUseCase } from './usecase/gemausecase';

export const gameState = new GameState();
const presenter = new GamePresenter(gameState);
export const usecase = new GameUseCase(presenter);
