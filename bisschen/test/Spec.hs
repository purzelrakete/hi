import Test.Tasty
import Test.Tasty.HUnit

import Lib

main = defaultMain $
  testGroup "Tests" [
    testGroup "Unit tests"
      [ testCase "Mine genesis block creates new block" $
          (mine Genesis "NOPE") @?= (Block 1 0 "NOPE" "HASH" Genesis)
      ]
  ]
