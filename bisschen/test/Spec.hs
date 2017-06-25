import Test.Tasty
import Test.Tasty.HUnit

import Lib

main = defaultMain $
  testGroup "Tests" [
    testGroup "Unit tests"
      [ testCase "Mine genesis block creates the second block" $
          (mine Genesis "tx") @?= (Block 1 0 "tx" "901131d838b17aac0f7885b81e03cbdc9f5157a00343d30ab22083685ed1416a" Genesis)
      ]
  ]
