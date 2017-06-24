module Lib
    ( start
    , mine
    , Block(Genesis)
    , Tx
    ) where

-- main
start :: IO ()
start = putStrLn "hello"

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

mine :: Block t -> Block t
mine Genesis = Genesis

-- ledger
data Tx = Tx
  { bisschens :: Integer
  , from :: Account
  , to :: Account
  } deriving (Show, Eq)
