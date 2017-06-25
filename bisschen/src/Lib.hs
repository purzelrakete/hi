module Lib
    ( start
    , mine
    , Block(Genesis, Block)
    , Tx
    ) where

import Crypto.Hash as H
import Crypto.Hash.Algorithms
import Data.ByteString.Char8 (pack)

-- main
start :: IO ()
start = Prelude.putStrLn "hello"

-- accounts
type Account = String
data Wallet = Wallet
  { account :: Account
  , privateKey :: String
  } deriving (Show)

-- blockchain
type Hash = String
data Block t = Genesis | Block
  { version :: Integer
  , nonce :: Integer
  , content :: t
  , hash :: Hash
  , previous :: Block t
  } deriving (Show, Eq)

mine :: Block t -> t -> Block t
mine Genesis content = Block 1 0 content (hashStr "GENESIS") Genesis

hashStr :: String -> Hash
hashStr x = show (H.hash (pack x) :: Digest SHA256)

-- ledger
data Tx = Tx
  { bisschens :: Integer
  , from :: Account
  , to :: Account
  } deriving (Show, Eq)
