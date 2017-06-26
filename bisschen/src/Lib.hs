module Lib
    ( mine
    , search
    , mkhash
    , zeros
    , Block(..)
    , Tx
    ) where

import Crypto.Hash.Algorithms (SHA256)
import qualified Data.ByteString.Char8 as C (pack)
import qualified Crypto.Hash as H
import qualified Data.ByteArray as BA

-- accounts
type Account = String
data Wallet = Wallet
  { account :: Account
  , key :: String
  } deriving Show

-- blockchain
type Hash = H.Digest SHA256
data Block t = Genesis | Block
  { version :: Int
  , nonce :: Int
  , content :: t
  , root :: Block t
  , digest :: Hash
  } deriving (Show, Eq)

mine :: Show t => Block t -> t -> Block t
mine Genesis x = Block
  { version = 1
  , nonce = n
  , content = x
  , root = Genesis
  , digest = d
  }
  where difficulty = 1
        (n, d) = search 0 difficulty

-- search for a hash
search :: Int -> Int -> (Int, Hash)
search x difficulty
  | difficulty < 0 || difficulty > 256 = error "invalid difficulty"
  | zeros hsh >= difficulty = (x, hsh)
  | otherwise = search next difficulty
  where next = x + 1
        hsh = mkhash x

-- number of leading zeros in a hash
zeros :: Hash -> Int
zeros hsh = f $ BA.unpack hsh
  where f [] = 0
        f xs | head xs ==  0 = (+2) . f . tail $ xs
             | head xs <= 15 = 1
             | otherwise = 0

-- ledger
data Tx = Tx
  { bisschens :: Int
  , from :: Account
  , to :: Account
  } deriving (Show, Eq)

-- hash
mkhash :: Int -> Hash
mkhash x = H.hash $ C.pack $ show x
