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
import { KOMA_LARGE_1, KOMA_LARGE_2, KOMA_MEDIUM_1, KOMA_MEDIUM_2, KOMA_SMALL_1, KOMA_SMALL_2 } from './lib/util/const';

function App() {
  const [turn, setTurn] = useState(1);
  const [winner, setWinner] = useState(0);
  const [, setBoard] = useState<BoardInfo[][]>([]);
  const [boardImg, setBoardImg] = useState<string[]>([]);
  const [selectKoma, setSelectKoma] = useState<string[]>([]);
  const [komaSize, setKomaSize] = useState(-1);

  const [playerImg1, setImg1] = useState('');
  const [playerImg2, setImg2] = useState('');
  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    (() => {
      usecase.init();
      setBoard(gameState.board);
      setBoardImg(gameState.boardImg);
      setImg1('/stamp.png');
      setImg2('/beer.png');
    })();
  }, []);

  const tapped = (index: number) => {
    if(komaSize === -1){
      alert("コマを選択してください");
      return;
    }
    const koma = usecase.input(index, turn,komaSize, gameState.board, gameState.boardImg);
    if (koma.turn == -1) return;

    setBoard(gameState.board);
    setBoardImg(gameState.boardImg);

    if (usecase.isWin(gameState.board)) {
      onOpen();
      setWinner(turn);
      reset();
      return;
    }

    if (turn === 1) {
      setTurn(2);
    } else if (turn === 2) {
      setTurn(1);
    }
    const targetId = selectKoma.findIndex((k) => k === '選択中');
    const updatedKoma =[...selectKoma];
    updatedKoma[targetId] = '済み';
    setSelectKoma(updatedKoma);
    setKomaSize(-1);
  };

  const reset = () => {
    setTurn(1);
    usecase.init();
    setBoard(gameState.board);
    setBoardImg(gameState.boardImg);
    setSelectKoma([]);
  };

  const closeModal = () => {
    onClose();
  };

  const select = (id: number, koma: string[]) => {
    const updatedKoma = [...koma];

    if(updatedKoma[id] === '' || updatedKoma[id] === undefined) {
      updatedKoma[id] = '選択中';  
    } else {
      updatedKoma[id] = '';
    }

    switch(id) {
      case 0:
        setKomaSize(KOMA_LARGE_1);
        break;
      case 1:
        setKomaSize(KOMA_LARGE_2);
        break;
      case 2:
        setKomaSize(KOMA_MEDIUM_1);
        break;
      case 3:
        setKomaSize(KOMA_MEDIUM_2);
        break;
      case 4:
        setKomaSize(KOMA_SMALL_1);
        break;
      case 5:
        setKomaSize(KOMA_SMALL_2);
        break;
    }
    setSelectKoma(updatedKoma);
  };

  return (
    <Box p={4}>
      <Heading mb={4} textAlign="center">
        五目並べ
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
                <img
                  src={boardImg[rowIndex] === 'red.200' ? playerImg1 : playerImg2} // 画像のURL
                  alt="Sample Image" // 画像の代替テキスト
                  width="100%" // 画像の幅（ボックスに合わせて100%にする）
                  height="100%" // 画像の高さ（ボックスに合わせて100%にする）
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
                w="80px"
                h="80px"
                border="1px solid #ccc"
                bg="teal.200"
                cursor="pointer"
                borderRadius="md"
                boxShadow="md"
                onClick={() => select(id, selectKoma)}
              >
                {selectKoma[id]}
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
