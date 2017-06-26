import Test.Tasty
import Test.Tasty.HUnit
import Test.Tasty.QuickCheck as QC
import Test.Tasty.SmallCheck as SC

import Data.List
import Data.Ord

import Lib

main = defaultMain $
  testGroup "Tests"
    [
      testGroup "Property tests"
       [ SC.testProperty "hash difficulty should be satisfied" $
           \nonce difficulty ->
             let gotZeros (_, h) = leadingZeros $ show h
                 leadingZeros xs = length $ takeWhile (\x -> x == '0') xs
              in (gotZeros $ search nonce difficulty) == max 0 difficulty
       ]
    , testGroup "Unit tests"
       [ testCase "Mine genesis block creates the second block" $
           mine Genesis "tx" @?= Block
             { version = 1
             , nonce = 39
             , content = "tx"
             , root = Genesis
             , digest = mkhash 39
             }
       ]
  ]
