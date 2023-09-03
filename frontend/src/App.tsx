import {
  Box,
  Button,
  Grid,
  Heading,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  useDisclosure
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import { gameState, usecase } from './lib/container';
import { BoardInfo } from './lib/domain/entity';
import {
  KOMA_LARGE_1,
  KOMA_LARGE_2,
  KOMA_MEDIUM_1,
  KOMA_MEDIUM_2,
  KOMA_SMALL_1,
  KOMA_SMALL_2,
  pch_List
} from './lib/util/const';

function App() {
  const [turn, setTurn] = useState(1);
  const [winner, setWinner] = useState(0);
  const [, setBoard] = useState<BoardInfo[][]>([]);
  const [boardImg, setBoardImg] = useState<string[]>([]);
  const [boardIcon, setBoardIcon] = useState<string[]>([]);
  const [selectKoma_P1, setSelectKomaP1] = useState<string[]>([]);
  const [selectKoma_P2, setSelectKomaP2] = useState<string[]>([]);
  const [komaSize, setKomaSize] = useState(-1);

  const [, setImg1] = useState('');
  const [, setImg2] = useState('');
  const { isOpen, onOpen, onClose } = useDisclosure();

  const title = 'ゴブレットゴブラーズ';

  useEffect(() => {
    (() => {
      usecase.init();
      setBoard(gameState.board);
      setBoardImg(gameState.boardImg);
      setBoardIcon(gameState.boardIcon);
      setImg1('/stamp.png');
      setImg2('/beer.png');
    })();
  }, []);

  const tapped = (index: number) => {
    if (komaSize === -1) {
      alert('コマを選択してください');
      return;
    }
    const koma = usecase.input(index, turn, komaSize, gameState.board, gameState.boardImg, gameState.boardIcon);
    if (koma.turn == -1) return;

    setBoard(gameState.board);
    setBoardImg(gameState.boardImg);
    setBoardIcon(gameState.boardIcon);

    if (usecase.isWin(gameState.board)) {
      onOpen();
      setWinner(turn);
      return;
    }

    if (turn === 1) {
      setTurn(2);
    } else if (turn === 2) {
      setTurn(1);
    }
    const targetId =
      turn === 1 ? selectKoma_P1.findIndex((k) => k === '選択中') : selectKoma_P2.findIndex((k) => k === '選択中');
    const updatedKoma = turn === 1 ? [...selectKoma_P1] : [...selectKoma_P2];
    updatedKoma[targetId] = '済み';
    turn === 1 ? setSelectKomaP1(updatedKoma) : setSelectKomaP2(updatedKoma);
    setKomaSize(-1);
  };

  const reset = () => {
    setTurn(1);
    usecase.init();
    setBoard(gameState.board);
    setBoardImg(gameState.boardImg);
    setSelectKomaP1([]);
    setSelectKomaP2([]);
  };

  const closeModal = () => {
    onClose();
  };

  const select = (id: number, koma: string[]) => {
    const updatedKoma = [...koma];
    const targetId =
      turn === 1 ? selectKoma_P1.findIndex((k) => k === '選択中') : selectKoma_P2.findIndex((k) => k === '選択中');

    if (targetId !== -1) {
      updatedKoma[targetId] = '';
    }

    if (updatedKoma[id] === '' || updatedKoma[id] === undefined) {
      updatedKoma[id] = '選択中';
    } else {
      updatedKoma[id] = '';
    }

    switch (id) {
      case 0:
        setKomaSize(KOMA_LARGE_1);
        turn === 1 ? setImg1('/Gophersvg_pink.svg') : setImg2('/Gophersvg_pink.svg');
        break;
      case 1:
        turn === 1 ? setImg1('/Gophersvg_pink.svg') : setImg2('/Gophersvg_pink.svg');
        setKomaSize(KOMA_LARGE_2);
        break;
      case 2:
        turn === 1 ? setImg1('/Gophersvg_yellow.svg') : setImg2('/Gophersvg_yellow.svg');
        setKomaSize(KOMA_MEDIUM_1);
        break;
      case 3:
        turn === 1 ? setImg1('/Gophersvg_yellow.svg') : setImg2('/Gophersvg_yellow.svg');
        setKomaSize(KOMA_MEDIUM_2);
        break;
      case 4:
        turn === 1 ? setImg1('/Gophersvg.svg') : setImg2('/Gophersvg.svg');
        setKomaSize(KOMA_SMALL_1);
        break;
      case 5:
        turn === 1 ? setImg1('/Gophersvg.svg') : setImg2('/Gophersvg.svg');
        setKomaSize(KOMA_SMALL_2);
        break;
    }
    turn === 1 ? setSelectKomaP1(updatedKoma) : setSelectKomaP2(updatedKoma);
  };

  return (
    <Box p={4}>
      <Heading mb={4} textAlign="center">
        {title}
      </Heading>
      <Box display="flex" flexDirection="column" justifyContent="center" alignItems="center">
        <h1>{turn} のターン</h1>
        <Grid templateColumns="repeat(3, 0fr)" gap={1} marginTop="30px">
          {[...Array(9)].map((_, rowIndex) => (
            <Box
              key={rowIndex}
              h="80px"
              w="80px"
              border="1px solid #ccc"
              display="flex"
              alignItems="center"
              justifyContent="center"
              cursor="pointer"
              bg={boardImg[rowIndex]}
              borderRadius="md"
              boxShadow="md"
              onClick={() => tapped(rowIndex)}
            >
              {boardImg[rowIndex] === '' ? (
                <></>
              ) : (
                // 盤面
                <img
                  src={boardIcon[rowIndex]} // 画像のURL
                  alt="Sample Image" // 画像の代替テキスト
                  width="80%" // 画像の幅（ボックスに合わせて100%にする）
                  height="80%" // 画像の高さ（ボックスに合わせて100%にする）
                />
              )}
            </Box>
          ))}
        </Grid>

        <Box w="700px" display="flex" justifyContent="center">
          <Grid templateColumns="repeat(6, 1fr)" gap={1} marginTop="30px">
            {[...Array(6)].map((_, id) => (
              <Box
                key={id}
                position="relative"
                w="80px"
                h="80px"
                onClick={() => (turn === 1 ? select(id, selectKoma_P1) : select(id, selectKoma_P2))}
              >
                <Box
                  position="absolute"
                  top="0"
                  left="0"
                  w="80px"
                  h="80px"
                  border="1px solid #ccc"
                  bg={turn === 1 ? 'red.200' : 'blue.200'}
                  cursor="pointer"
                  borderRadius="md"
                  boxShadow="md"
                  textAlign="center"
                  zIndex={0}
                />

                <Box position="absolute" top="0" left="0" w="100%" h="80%" zIndex={1}>
                  {turn === 1 && selectKoma_P1[id] === '済み' ? (
                    <></>
                  ) : turn === 2 && selectKoma_P2[id] === '済み' ? (
                    <></>
                  ) : (
                    <img
                      // TODO: ここを変更する
                      // 選択肢
                      src={pch_List[id]} // 画像のURL
                      alt="Sample Image" // 画像の代替テキスト
                      width="80%" // 画像の幅（ボックスに合わせて100%にする）
                      height="80%" // 画像の高さ（ボックスに合わせて100%にする）
                    />
                  )}
                </Box>
              </Box>
            ))}
            {[...Array(6)].map((_, id) => (
              <Box
                key={id}
                h="40px"
                w="80px"
                border="1px solid #ccc"
                bg="teal.200"
                fontWeight="bold"
                textAlign="center"
              >
                {turn === 1 ? selectKoma_P1[id] : selectKoma_P2[id]}
              </Box>
            ))}
          </Grid>
        </Box>

        <Box w="700px" display="flex" justifyContent="flex-end">
          <Button marginTop="30px" bg="orange.200" onClick={() => reset()}>
            リセット
          </Button>
        </Box>

        <Modal isOpen={isOpen} onClose={onClose}>
          <ModalOverlay />
          <ModalContent>
            <ModalHeader>ゲーム結果</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
              <Stack spacing={3}>
                <h2>🎉🎉🎉 {winner === 1 ? '先攻' : '後攻'}の勝利！！！🎉🎉🎉</h2>
              </Stack>
            </ModalBody>

            <ModalFooter>
              <Button variant="ghost" onClick={() => closeModal()}>
                閉じる
              </Button>
            </ModalFooter>
          </ModalContent>
        </Modal>
      </Box>
    </Box>
  );
}

export default App;
