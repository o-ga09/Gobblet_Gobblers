import {
  Box,
  Button,
  Grid,
  GridItem,
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

function App() {
  const [turn, setTurn] = useState(1);
  const [winner, setWinner] = useState(0);
  const [, setBoard] = useState<number[][]>([]);
  const [boardImg, setBoardImg] = useState<string[]>([]);
  const [koma, setKoma] = useState<string[]>([]);

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
    const koma = usecase.input(index, turn, gameState.board, gameState.boardImg);
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
  };

  const reset = () => {
    setTurn(1);
    usecase.init();
    setBoard(gameState.board);
    setBoardImg(gameState.boardImg);
    setKoma([]);
  };

  const closeModal = () => {
    onClose();
  };

  const selectKoma = (id: number, koma: string[]) => {
    console.log(koma);
    koma[id] = '選択中';
    const newKoma = koma;
    setKoma(newKoma);
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
                onClick={() => selectKoma(id, koma)}
              >
                {koma[id]}
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
