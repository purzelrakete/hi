module Main where

import Test.Tasty
import Test.Tasty.HUnit

import Lib

main :: IO ()
main = do
  defaultMain (testGroup "Blockchain" [mineGenesisBlockTest])

mineGenesisBlockTest :: TestTree
mineGenesisBlockTest = testCase "Testing mining blocks"
  (assertEqual "should mine genesis block" (Genesis::Block ()) (mine Genesis))
